package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/shawnpeng17/17vpn/internal/pritunl"
)

var rootCmd = &cobra.Command{
	Use:   "17vpn",
	Short: "17vpn tool",
	Run: func(cmd *cobra.Command, args []string) {
		p := pritunl.New()
		profiles := p.Profiles()
		conns := p.Connections()

		if err := list(profiles, conns); err != nil {
			color.Yellow(err.Error())
			return
		}

		var id string
		_, _ = color.New(color.FgYellow).Print("Enter ID or Server: ")
		_, _ = fmt.Scanln(&id)

		if id == "" {
			return
		}

		// check profile exist
		var targetProfile pritunl.Profile
		isActionDisconnect := false
		for i, profile := range profiles {
			if strconv.Itoa(i+1) == id || strings.ToUpper(id) == profile.Server {
				targetProfile = profile
				isActionDisconnect = conns[profile.ID].Status == "connected"
				break
			}
		}
		if targetProfile == (pritunl.Profile{}) {
			color.Red("Profile not exists!")
			return
		}

		if isActionDisconnect {
			color.White("Disconnecting %s...", targetProfile.Server)
			p.Disconnect(targetProfile.ID)
			return
		}

		// disconnect all connection
		for _, profile := range profiles {
			if _, ok := conns[profile.ID]; ok {
				color.White("Disconnecting %s...", profile.Server)
				p.Disconnect(profile.ID)
			}
		}

		// connect target profile
		color.Yellow("Connecting %s...", targetProfile.Server)
		p.Connect(targetProfile.ID, "password")

		timeout := time.NewTimer(time.Minute)

	Loop:
		for {
			select {
			case <-timeout.C:
				color.Red("Connect %s timeout!", targetProfile.Server)
				break Loop
			default:
				status := p.Connections()[targetProfile.ID].Status
				switch status {
				case "connected":
					color.Green("Connect %s completed!", targetProfile.Server)
					break Loop
				case "":
					color.Red("Connect %s failed!", targetProfile.Server)
					break Loop
				}
				time.Sleep(time.Second)
			}
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
