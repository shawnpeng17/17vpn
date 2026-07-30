package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/shawnpeng17/17vpn/internal/pritunl"
	"github.com/shawnpeng17/17vpn/internal/tui"
)

var rootCmd = &cobra.Command{
	Use:   "17vpn",
	Short: "17vpn tool",
	Run: func(cmd *cobra.Command, args []string) {
		if err := initConfig(); err != nil {
			color.Red(err.Error())
			return
		}

		p := pritunl.New()
		if err := tui.Run(p, password); err != nil {
			color.Red(err.Error())
			return
		}
	},
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		return
	}
}
