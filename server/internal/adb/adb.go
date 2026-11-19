package adb

import (
	"flag"
	"strings"

	"github.com/Teages/vTablet/internal/logger"
	. "github.com/Teages/vTablet/internal/tools"
	adb "github.com/zach-klippenstein/goadb"
)

const (
	AppPackage  = "xyz.teages.vtablet"
	AppActivity = "xyz.teages.vtablet/.MainActivity"
	ClientUrl   = "https://vtablet.yoruakio.xyz?server=localhost:23052"
)

var (
	adbClient *adb.Adb
)

func Watch(f func(adb.DeviceStateChangedEvent)) {
	if adbClient == nil {
		logger.Warn("Can't watch: adb server not found")
		return
	}
	go func() {
		for {
			w := adbClient.NewDeviceWatcher()
			for event := range w.C() {
				// logger.Log("\t[%s]%+v\n", time.Now(), event)
				f(event)
			}
		}
	}()
}

func AutoConnect(c func(string)) {
	Watch(func(event adb.DeviceStateChangedEvent) {
		if event.NewState == adb.StateOnline {
			logger.Log(event.Serial)
			c(event.Serial)
		}
	})
}

func Exec(args ...string) {
	RunCommand(adbPath, args...)
}

func ExecOutput(args ...string) (string, error) {
	return RunCommandOutput(adbPath, args...)
}

func init() {
	Exec("start-server")
	a, err := adb.NewWithConfig(adb.ServerConfig{
		PathToAdb: adbPath,
		Port:      *flag.Int("p", adb.AdbPort, ""),
	})
	if logger.Catch(err) {
		logger.Error("Can't start ADB server")
		adbClient = nil
		return
	}
	adbClient = a
	AutoConnect(func(serial string) {
		Exec("-s", serial, "reverse", "tcp:23052", "tcp:23052")
		// Check if the vTablet native client app is installed
		out, err := ExecOutput("-s", serial, "shell", "pm", "path", AppPackage)
		if err == nil && strings.Contains(out, "package:") {
			logger.Log("Client app detected on %s, launching native app...", serial)
			Exec("-s", serial, "shell", "am", "start", "-n", AppActivity, "-a", "android.intent.action.VIEW", "-d", ClientUrl)
		} else {
			logger.Log("Client app not found on %s, opening in browser...", serial)
			Exec("-s", serial, "shell", "am", "start", "-a", "android.intent.action.VIEW", "-d", ClientUrl)
		}
	})
}

func Restart() {
	if adbClient == nil {
		logger.Warn("Can't restart ADB: adb server not found")
		return
	}
	adbClient.KillServer()
	Exec("start-server")
}
