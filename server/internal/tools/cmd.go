package tools

import (
	"os"

	"github.com/Teages/vTablet/internal/logger"
)

func GetSelfPath() string {
	exePath, err := os.Executable()
	logger.Catch(err)
	return exePath
}
