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
	var finalArgs []string
	for _, arg := range args {
		parts := strings.Fields(arg)
		finalArgs = append(finalArgs, parts...)
	}
	e := exec.Command(exePath, finalArgs...)
	logger.Log("Running: %s %s", exePath, strings.Join(finalArgs, " "))
	e.Stdout = os.Stdout
	e.Stderr = os.Stderr
	err := e.Run()
	logger.Catch(err)
}
