package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Global *Config
)

type Config struct {
	DBConfig               DBConfig               `mapstructure:"db"`
	LogLevel               string                 `mapstructure:"log_level"`
	Port                   int                    `mapstructure:"port"`
	GoogleDriveProxyConfig GoogleDriveProxyConfig `mapstructure:"google"`

	// dynamic configs
	BaseApiUrl string `mapstructure:"-"`
}

func Init(cfgFile string, rootCmd cobra.Command) {
	InitViper(cfgFile, rootCmd)

	if Global == nil {
		Global = &Config{}
	} else {
		panic(fmt.Errorf("configs already initialized"))
	}

	// load config into Configs
	if err := viper.Unmarshal(Global); err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}

	initConfigFields()
}

func initConfigFields() {
	if Global == nil {
		panic(fmt.Errorf("configs not initialized"))
	}

	if Global.DBConfig.Filename == "" {
		Global.DBConfig.Filename = "lfsproxy.db"
	}

	if Global.DBConfig.Username == "" {
		Global.DBConfig.Username = "admin"
	}

	if Global.DBConfig.Password == "" {
		Global.DBConfig.Password = "password"
	}

	if Global.DBConfig.EncryptionAlgorithm == "" {
		Global.DBConfig.EncryptionAlgorithm = "SHA256"
	}

	if Global.DBConfig.Timeout == 0 {
		Global.DBConfig.Timeout = 5
	}

	// root api url
	Global.BaseApiUrl = fmt.Sprintf("http://localhost:%d", Global.Port)
}
