package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/Teages/go-vfile"
	_ "github.com/Teages/vTablet/internal/adb"
	"github.com/Teages/vTablet/internal/logger"
)

var (
	headless = flag.Bool("headless", false, "Run in headless mode without system tray")
	verbose  = flag.Bool("verbose", false, "Enable verbose logging for pointer packets")
)

func main() {
	flag.Parse()

	go initWebServices(23052)

	if *headless {
		logger.Log("vTablet server running in headless mode")
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		logger.Log("Shutting down vTablet server...")
		vfile.Close()
		return
	}

	initTray(func() {
		// onExit
		vfile.Close()
	})
}
