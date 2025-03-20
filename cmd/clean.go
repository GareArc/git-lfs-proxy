package cmd

import (
	"os"

	"github.com/GareArc/git-lfs-proxy/config"
	"github.com/spf13/cobra"
)

var (
	cleanCmd = &cobra.Command{
		Use:   "clean",
		Short: "Clean the cache",
		Long:  "Clean the cache",
		Run: func(cmd *cobra.Command, args []string) {
			os.RemoveAll(config.CONFIG_DIR)
			os.MkdirAll(config.CONFIG_DIR, 0755)
		},
	}
)
