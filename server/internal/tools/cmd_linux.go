//go:build !windows

package tools

import (
	"os"
	"os/exec"
	"strings"

	"github.com/Teages/vTablet/internal/logger"
)

func RunCommand(path string, args ...string) {
	exePath := strings.TrimSpace(path)
	e := exec.Command(exePath, args...)
	logger.Log("Running: %s %s", exePath, strings.Join(args, " "))
	e.Stdout = os.Stdout
	e.Stderr = os.Stderr
	err := e.Run()
	logger.Catch(err)
}
