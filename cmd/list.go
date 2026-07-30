package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/shawnpeng17/17vpn/internal/pritunl"
)

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "List profiles",
	Run: func(cmd *cobra.Command, args []string) {
		p := pritunl.New()

		conns, err := p.Connections()
		if err != nil {
			color.Yellow(err.Error())
			return
		}
		if err := list(p.Profiles(), conns); err != nil {
			color.Yellow(err.Error())
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
