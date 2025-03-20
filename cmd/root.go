package cmd

import (
	"os"
	"os/signal"
	"runtime"

	"github.com/GareArc/git-lfs-proxy/config"
	"github.com/GareArc/git-lfs-proxy/internal/db"
	"github.com/GareArc/git-lfs-proxy/internal/logging"
	"github.com/GareArc/git-lfs-proxy/internal/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"net/http"
	_ "net/http/pprof"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "lfsproxy",
		Short: "Start LFS Proxy Server",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			// init logging
			logging.Init(config.Global.LogLevel)

			// ppfrof
			if config.Global.LogLevel == "debug" {
				runtime.SetBlockProfileRate(1)
				runtime.SetMutexProfileFraction(1)
				go func() {
					log.Info().Msg("Starting pprof server...")
					if err := http.ListenAndServe("localhost:6060", nil); err != nil {
						log.Error().Err(err).Msg("pprof server error")
					}
				}()
			}

			// init db
			db.InitDB()

			log.Info().Msg("Starting LFS Proxy Server...")
			server, cancel := server.NewServer()
			defer cancel()

			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt)

			go func() {
				<-sigChan
				cancel()
			}()

			// run http server
			if err := server.Start(); err != nil {
				log.Error().Err(err).Msg("Server Error")
				cancel()
			}

			log.Info().Msg("LFS Proxy Server down")
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/lfsproxy.config.toml)")
	rootCmd.PersistentFlags().StringP("log-level", "l", "", "log level (debug, info, warn, error, fatal, panic)")
	rootCmd.PersistentFlags().IntP("port", "p", 8080, "port to listen on")

	// add subcommands
	rootCmd.AddCommand(cleanCmd)
}

func initConfig() {
	config.Init(cfgFile, *rootCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
