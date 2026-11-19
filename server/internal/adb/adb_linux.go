//go:build !windows

package adb

import (
	"os/exec"
)

var (
	adbPath = func() string {
		p, err := exec.LookPath("adb")
		if err == nil {
			return p
		}
		return "adb"
	}()
)
