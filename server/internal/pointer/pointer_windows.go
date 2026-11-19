//go:build windows

package pointer

import (
	"fmt"

	"github.com/Teages/go-vdigi"
)

type windowsBackend struct {
	screenId int
	device   *vdigi.Pointer
}

func GetScreens() []ScreenData {
	var screens []ScreenData
	count := vdigi.GetScreens().GetScreenCount()
	for i := 0; i < count; i++ {
		s, _ := vdigi.GetScreens().GetScreen(i)
		screens = append(screens, ScreenData{
			Uid:    s.Uid,
			Path:   fmt.Sprintf("/digi/%s", s.Uid),
			Width:  s.Width,
			Height: s.Height,
		})
	}
	return screens
}

func Create(screenUid string) *AbsPointer {
	screenId, _ := vdigi.GetScreens().GetScreenIdByUid(screenUid)
	device, _ := vdigi.CreatePointerForScreen(screenId)
	return &AbsPointer{
		screenUid: screenUid,
		backend: &windowsBackend{
			screenId: screenId,
			device:   device,
		},
	}
}

func (b *windowsBackend) Update(X, Y int32, Pressure uint32, TiltX, TiltY int32) error {
	screen, _ := vdigi.GetScreens().GetScreen(b.screenId)
	width := int32(screen.Width)
	height := int32(screen.Height)
	x := X * width / 32767
	y := Y * height / 32767
	tiltX := int32(float64(TiltX) * 90.0 / 32767.0)
	tiltY := int32(float64(TiltY) * 90.0 / 32767.0)
	return b.device.UpdateWithTilt(x, y, Pressure, tiltX, tiltY)
}

func (b *windowsBackend) Destroy() {
	if b.device != nil {
		b.device.Destroy()
	}
}
