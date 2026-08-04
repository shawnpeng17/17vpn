package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of 17vpn",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("17vpn version %s (commit: %s, built at: %s)\n", version, commit, date)
	},
}

func init() {
	rootCmd.Version = version
	rootCmd.AddCommand(versionCmd)
}
