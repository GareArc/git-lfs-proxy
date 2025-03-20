package config

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func InitViper(cfgFile string, rootCmd cobra.Command) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// use config directory if config file is not given
		viper.AddConfigPath(CONFIG_DIR)
		viper.SetConfigName(".lfsproxy.config")
		cfgFile = fmt.Sprintf("%s/.lfsproxy.config.toml", CONFIG_DIR)
	}
	viper.SetConfigType("toml")
	viper.AutomaticEnv()
	initViperFields()

	if err := viper.ReadInConfig(); err != nil {
		// create config file if not exists
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found, creating one")
			// create the config directory if not exists
			if err := os.MkdirAll(CONFIG_DIR, 0755); err != nil {
				panic(fmt.Errorf("failed to create config directory: %v", err))
			}

			if err := viper.WriteConfigAs(cfgFile); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	bindFlags(rootCmd)

}

func bindFlags(rootCmd cobra.Command) {
	// Bind the flags to the viper
	cmdLogLevel := rootCmd.PersistentFlags().Lookup("log-level")
	if cmdLogLevel.Value.String() != "" {
		viper.BindPFlag("log_level", cmdLogLevel)
	}
	cmdPortLevel := rootCmd.PersistentFlags().Lookup("port")
	if cmdPortLevel.Value.String() != "" {
		viper.BindPFlag("port", cmdPortLevel)
	}
}

func initViperFields() {
	viper.SetDefault("log_level", "info")
	viper.SetDefault("port", 8080)
	viper.SetDefault("google.enabled", false)
}
