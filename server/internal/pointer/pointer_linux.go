//go:build !windows

package pointer

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"github.com/Teages/vTablet/internal/logger"
)

const (
	UI_SET_EVBIT   = 0x40045564
	UI_SET_KEYBIT  = 0x40045565
	UI_SET_ABSBIT  = 0x40045567
	UI_SET_PROPBIT = 0x4004556e
	UI_DEV_SETUP   = 0x405c5503
	UI_ABS_SETUP   = 0x401c5504
	UI_DEV_CREATE  = 0x5501
	UI_DEV_DESTROY = 0x5502

	EV_SYN = 0x00
	EV_KEY = 0x01
	EV_ABS = 0x03

	SYN_REPORT = 0x00

	BTN_TOOL_PEN = 0x140
	BTN_TOUCH    = 0x14a
	BTN_STYLUS   = 0x14b
	BTN_STYLUS2  = 0x14c

	ABS_X        = 0x00
	ABS_Y        = 0x01
	ABS_PRESSURE = 0x18
	ABS_TILT_X   = 0x1a
	ABS_TILT_Y   = 0x1b

	INPUT_PROP_DIRECT = 0x01
	BUS_VIRTUAL       = 0x06

	DeviceName = "vtablet-digitizer"
)

type uinputSetup struct {
	Bustype      uint16
	Vendor       uint16
	Product      uint16
	Version      uint16
	Name         [80]byte
	FfEffectsMax uint32
}

type uinputAbsSetup struct {
	Code    uint16
	_       [2]byte
	Value   int32
	Minimum int32
	Maximum int32
	Fuzz    int32
	Flat    int32
	Res     int32
}

type inputEvent struct {
	Time  [16]byte
	Type  uint16
	Code  uint16
	Value int32
}

func ioctl(fd int, req uintptr, arg uintptr) error {
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), req, arg)
	if err != 0 {
		return err
	}
	return nil
}

type linuxUinputBackend struct {
	mu     sync.Mutex
	file   *os.File
	screen ScreenData
	closed bool
}

func GetScreens() []ScreenData {
	// 1. Try wlr-randr --json (standard for wlroots/Wayland: Hyprland, Sway, Mango, etc.)
	if out, err := exec.Command("wlr-randr", "--json").Output(); err == nil {
		var list []struct {
			Name    string `json:"name"`
			Enabled bool   `json:"enabled"`
			Modes   []struct {
				Width   int  `json:"width"`
				Height  int  `json:"height"`
				Current bool `json:"current"`
			} `json:"modes"`
		}
		if json.Unmarshal(out, &list) == nil {
			var screens []ScreenData
			for _, m := range list {
				if !m.Enabled {
					continue
				}
				w, h := 0, 0
				for _, mode := range m.Modes {
					if mode.Current {
						w, h = mode.Width, mode.Height
						break
					}
				}
				if w > 0 && h > 0 {
					uid := fmt.Sprintf("%03x%03x", w, h)
					screens = append(screens, ScreenData{
						Uid:    uid,
						Path:   fmt.Sprintf("/digi/%s", uid),
						Width:  w,
						Height: h,
						Name:   m.Name,
					})
				}
			}
			if len(screens) > 0 {
				return screens
			}
		}
	}

	// 2. Try hyprctl monitors -j
	if out, err := exec.Command("hyprctl", "monitors", "-j").Output(); err == nil {
		var list []struct {
			Name   string `json:"name"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		}
		if json.Unmarshal(out, &list) == nil {
			var screens []ScreenData
			for _, m := range list {
				if m.Width > 0 && m.Height > 0 {
					uid := fmt.Sprintf("%03x%03x", m.Width, m.Height)
					screens = append(screens, ScreenData{
						Uid:    uid,
						Path:   fmt.Sprintf("/digi/%s", uid),
						Width:  m.Width,
						Height: m.Height,
						Name:   m.Name,
					})
				}
			}
			if len(screens) > 0 {
				return screens
			}
		}
	}

	// 3. Fallback: check /sys/class/drm/*/status + modes
	statusFiles, _ := filepath.Glob("/sys/class/drm/*/status")
	var screens []ScreenData
	for _, sf := range statusFiles {
		data, err := os.ReadFile(sf)
		if err == nil && strings.TrimSpace(string(data)) == "connected" {
			dir := filepath.Dir(sf)
			modesFile := filepath.Join(dir, "modes")
			modeData, err := os.ReadFile(modesFile)
			if err == nil {
				lines := strings.Split(strings.TrimSpace(string(modeData)), "\n")
				if len(lines) > 0 && strings.Contains(lines[0], "x") {
					parts := strings.Split(lines[0], "x")
					if len(parts) == 2 {
						w, _ := strconv.Atoi(parts[0])
						h, _ := strconv.Atoi(parts[1])
						if w > 0 && h > 0 {
							name := filepath.Base(dir)
							uid := fmt.Sprintf("%03x%03x", w, h)
							screens = append(screens, ScreenData{
								Uid:    uid,
								Path:   fmt.Sprintf("/digi/%s", uid),
								Width:  w,
								Height: h,
								Name:   name,
							})
						}
					}
				}
			}
		}
	}
	if len(screens) > 0 {
		return screens
	}

	// 4. Default fallback: 1920x1080
	return []ScreenData{
		{
			Uid:    "780438", // 1920x1080 in hex
			Path:   "/digi/780438",
			Width:  1920,
			Height: 1080,
			Name:   "Default",
		},
	}
}

func Create(screenUid string) *AbsPointer {
	screens := GetScreens()
	targetScreen := ScreenData{
		Uid:    screenUid,
		Path:   fmt.Sprintf("/digi/%s", screenUid),
		Width:  1920,
		Height: 1080,
	}
	for _, s := range screens {
		if s.Uid == screenUid {
			targetScreen = s
			break
		}
	}

	f, err := os.OpenFile("/dev/uinput", os.O_WRONLY|syscall.O_NONBLOCK, 0660)
	if err != nil {
		logger.Error("Failed to open /dev/uinput: %v", err)
		return &AbsPointer{
			screenUid: screenUid,
			backend:   nil,
		}
	}

	fd := int(f.Fd())

	// Enable event types
	ioctl(fd, UI_SET_EVBIT, EV_SYN)
	ioctl(fd, UI_SET_EVBIT, EV_KEY)
	ioctl(fd, UI_SET_EVBIT, EV_ABS)

	// Enable keys
	ioctl(fd, UI_SET_KEYBIT, BTN_TOOL_PEN)
	ioctl(fd, UI_SET_KEYBIT, BTN_TOUCH)
	ioctl(fd, UI_SET_KEYBIT, BTN_STYLUS)
	ioctl(fd, UI_SET_KEYBIT, BTN_STYLUS2)

	// Set direct property
	ioctl(fd, UI_SET_PROPBIT, INPUT_PROP_DIRECT)

	// Set ABS axes
	absAxes := []struct {
		code uint16
		min  int32
		max  int32
	}{
		{ABS_X, 0, 32767},
		{ABS_Y, 0, 32767},
		{ABS_PRESSURE, 0, 8192},
		{ABS_TILT_X, -90, 90},
		{ABS_TILT_Y, -90, 90},
	}

	for _, a := range absAxes {
		ioctl(fd, UI_SET_ABSBIT, uintptr(a.code))
		s := uinputAbsSetup{
			Code:    a.code,
			Minimum: a.min,
			Maximum: a.max,
		}
		var buf bytes.Buffer
		binary.Write(&buf, binary.LittleEndian, s)
		ioctl(fd, UI_ABS_SETUP, uintptr(unsafe.Pointer(&buf.Bytes()[0])))
	}

	var devSetup uinputSetup
	devSetup.Bustype = BUS_VIRTUAL
	devSetup.Vendor = 0x1209
	devSetup.Product = 0x7879
	devSetup.Version = 1
	copy(devSetup.Name[:], DeviceName)

	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, devSetup)
	if err := ioctl(fd, UI_DEV_SETUP, uintptr(unsafe.Pointer(&buf.Bytes()[0]))); err != nil {
		logger.Error("UI_DEV_SETUP failed: %v", err)
		f.Close()
		return &AbsPointer{screenUid: screenUid, backend: nil}
	}

	if err := ioctl(fd, UI_DEV_CREATE, 0); err != nil {
		logger.Error("UI_DEV_CREATE failed: %v", err)
		f.Close()
		return &AbsPointer{screenUid: screenUid, backend: nil}
	}

	logger.Log("Created uinput virtual digitizer: %s for screen %s (%s)", DeviceName, targetScreen.Uid, targetScreen.Name)

	// Hyprland output mapping (non-fatal)
	if targetScreen.Name != "" {
		exec.Command("hyprctl", "keyword", fmt.Sprintf("device[%s]:output", DeviceName), targetScreen.Name).Run()
	}

	return &AbsPointer{
		screenUid: screenUid,
		backend: &linuxUinputBackend{
			file:   f,
			screen: targetScreen,
		},
	}
}

func (b *linuxUinputBackend) emit(evType uint16, code uint16, value int32) {
	if b.file == nil {
		return
	}
	ev := inputEvent{
		Type:  evType,
		Code:  code,
		Value: value,
	}
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, ev)
	b.file.Write(buf.Bytes())
}

func (b *linuxUinputBackend) Update(x, y int32, pressure uint32, tiltX, tiltY int32) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed || b.file == nil {
		return nil
	}

	// Clamp coordinates
	if x < 0 {
		x = 0
	} else if x > 32767 {
		x = 32767
	}
	if y < 0 {
		y = 0
	} else if y > 32767 {
		y = 32767
	}

	// Tool in-range
	b.emit(EV_KEY, BTN_TOOL_PEN, 1)

	if pressure > 0 {
		b.emit(EV_KEY, BTN_TOUCH, 1)
		b.emit(EV_ABS, ABS_PRESSURE, int32(pressure))
	} else {
		b.emit(EV_KEY, BTN_TOUCH, 0)
		b.emit(EV_ABS, ABS_PRESSURE, 0)
	}

	b.emit(EV_ABS, ABS_X, x)
	b.emit(EV_ABS, ABS_Y, y)

	tX := int32(float64(tiltX) * 90.0 / 32767.0)
	tY := int32(float64(tiltY) * 90.0 / 32767.0)
	b.emit(EV_ABS, ABS_TILT_X, tX)
	b.emit(EV_ABS, ABS_TILT_Y, tY)

	b.emit(EV_SYN, SYN_REPORT, 0)
	return nil
}

func (b *linuxUinputBackend) Destroy() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed || b.file == nil {
		return
	}
	b.closed = true

	// Lift pen
	b.emit(EV_KEY, BTN_TOUCH, 0)
	b.emit(EV_KEY, BTN_TOOL_PEN, 0)
	b.emit(EV_SYN, SYN_REPORT, 0)

	fd := int(b.file.Fd())
	ioctl(fd, UI_DEV_DESTROY, 0)
	b.file.Close()
	b.file = nil
	logger.Log("Destroyed uinput virtual digitizer")
}
