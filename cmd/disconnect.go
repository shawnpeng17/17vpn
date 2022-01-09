package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/shawnpeng17/17vpn/internal/pritunl"
)

var disconnectCmd = &cobra.Command{
	Use:     "disconnectAll",
	Aliases: []string{"d"},
	Short:   "Disconnect all connections",
	Run: func(cmd *cobra.Command, args []string) {
		p := pritunl.New()
		conns := p.Connections()
		if len(conns) == 0 {
			color.Yellow("No connection found!")
			os.Exit(1)
		}

		p.DisconnectAll()
		color.Green("All connections disconnected!")
	},
}

func init() {
	rootCmd.AddCommand(disconnectCmd)
}
