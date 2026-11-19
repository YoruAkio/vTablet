//go:build linux && !cgo

package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/Teages/vTablet/internal/logger"
)

var (
	appName    = "vTablet"
	appVersion = "v3.0.0 (Preview 4)"
)

func initTray(onExit func()) {
	logger.Log("vTablet running in CLI mode (built with CGO_ENABLED=0)")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	logger.Log("Shutting down...")
	onExit()
}

func updateConnectState(clientCount int) {
	// No-op without GUI tray
}

func openBrowser(url string) {
	err := exec.Command("xdg-open", url).Start()
	if err != nil {
		log.Printf("Failed to open browser: %v", err)
	}
}
