package cmd

import (
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/shawnpeng17/17vpn/internal/pritunl"
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
		profiles := p.Profiles()
		conns := p.Connections()

		if err := list(profiles, conns); err != nil {
			color.Yellow(err.Error())
			return
		}

		var options []string
		for _, profile := range profiles {
			options = append(options, profile.Server)
		}
		var id string
		prompt := &survey.Select{
			Message:       "Choose a Server",
			Options:       options,
			Default:       "",
			VimMode:       true,
		}
		if err := survey.AskOne(prompt, &id); err != nil {
			color.Red(err.Error())
			return
		}

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
				time.Sleep(time.Second)
			}
		}


		// connect target profile
		color.Yellow("Connecting %s...", targetProfile.Server)
		p.Connect(targetProfile.ID, password())

		timeout := time.NewTimer(30 * time.Second)

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
				time.Sleep(500 * time.Millisecond)
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
