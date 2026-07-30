package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/shawnpeng17/17vpn/internal/pritunl"
)

var disconnectCmd = &cobra.Command{
	Use:   "d",
	Short: "Disconnect all connections",
	Run: func(cmd *cobra.Command, args []string) {
		p := pritunl.New()
		conns, err := p.Connections()
		if err != nil {
			color.Yellow(err.Error())
			return
		}
		if len(conns) == 0 {
			color.Yellow("No connection found!")
			return
		}

		if err := p.DisconnectAll(); err != nil {
			color.Yellow(err.Error())
			return
		}
		color.Green("All connections disconnected!")
	},
}

func init() {
	rootCmd.AddCommand(disconnectCmd)
}
