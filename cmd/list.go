package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/shawnpeng17/17vpn/internal/pritunl"
)

var listCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List profiles",
	Run: func(cmd *cobra.Command, args []string) {
		p := pritunl.New()

		if err := list(p.Profiles(), p.Connections()); err != nil {
			color.Yellow(err.Error())
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
