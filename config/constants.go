package config

import (
	"fmt"
	"os"
)

var (
	CONFIG_DIR   string
	ClientID     string
	ClientSecret string
)

func init() {
	// set constants
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		panic(fmt.Errorf("failed to get user config directory: %v", err))
	}
	CONFIG_DIR = fmt.Sprintf("%s/git-lfs-proxy", userConfigDir)
}
