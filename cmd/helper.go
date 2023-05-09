package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
	"github.com/xlzd/gotp"

	"github.com/shawnpeng17/17vpn/internal/pritunl"
)

const (
	configFileName = ".17vpn.yaml"
)

type answer struct {
	Key string
	Pin string
}

func initConfig() error {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, configFileName)
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err == nil {
		return nil
	}

	var qs = []*survey.Question{
		{
			Name:   "key",
			Prompt: &survey.Password{Message: "Enter OTP key"},
		},
		{
			Name:   "pin",
			Prompt: &survey.Password{Message: "Enter Pin"},
		},
	}
	var ans answer
	if err := survey.Ask(qs, &ans); err != nil {
		return err
	}

	viper.Set("key", ans.Key)
	viper.Set("pin", ans.Pin)
	viper.SetConfigFile(path)
	if err := viper.WriteConfig(); err != nil {
		return err
	}

	color.Yellow("Config saved to %s\n\n", viper.ConfigFileUsed())
	return nil
}

func password(mode string) string {
	password := viper.GetString("pin")
	if mode == "otp_pin" {
		totp := gotp.NewDefaultTOTP(viper.GetString("key"))
		password += totp.Now()
	}
	return password
}

func list(profiles []pritunl.Profile, conns map[string]pritunl.Connection) error {
	if len(profiles) == 0 {
		return errors.New("No profile found in Pritunl!")
	}

	var rows [][]string
	for i, profile := range profiles {
		row := []string{strconv.Itoa(i + 1), profile.Server, profile.User, "disconnected", "", "", ""}
		if conn, ok := conns[profile.ID]; ok {
			row[3] = conn.Status
			if conn.Timestamp > 0 {
				row[4] = formatDuration(time.Since(time.Unix(conn.Timestamp, 0)))
			}
			row[5] = conn.ClientAddr
			row[6] = conn.ServerAddr
		}
		row[3] = formatStatus(row[3])
		rows = append(rows, row)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Server", "User", "Status", "Connected", "Client IP", "Server IP"})
	table.AppendBulk(rows)
	table.SetColWidth(1000)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetCenterSeparator(" ")
	table.SetColumnSeparator(" ")
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderLine(false)
	table.Render()

	return nil
}

func formatStatus(status string) string {
	status = strings.ToUpper(status)
	switch status {
	case "CONNECTED":
		return color.New(color.FgGreen, color.Bold).SprintfFunc()(status)
	case "CONNECTING":
		return color.New(color.FgYellow, color.Bold).SprintfFunc()(status + "...")
	case "DISCONNECTING":
		return color.New(color.FgBlack, color.Bold).SprintfFunc()(status + "...")
	case "DISCONNECTED":
		return color.New(color.FgBlack, color.Bold).SprintfFunc()(status)
	default:
		return color.New(color.FgBlack, color.Bold).SprintfFunc()("UNKNOWN")
	}
}

func formatDuration(sec time.Duration) string {
	var ret string
	day := 24 * time.Hour
	d := sec / day
	sec = sec % day
	h := sec / time.Hour
	sec = sec % time.Hour
	m := sec / time.Minute
	sec = sec % time.Minute
	s := sec / time.Second
	if d > 0 {
		ret += fmt.Sprintf("%dd", d)
	}
	if h > 0 {
		ret += fmt.Sprintf("%dh", h)
	}
	if m > 0 {
		ret += fmt.Sprintf("%dm", m)
	}
	ret += fmt.Sprintf("%ds", s)
	return ret
}
