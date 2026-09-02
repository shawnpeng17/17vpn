package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// applyBuildInfo fills in version, commit, and date from the Go build info
// when they're still at their defaults, since goreleaser's ldflags aren't
// applied when installing via `go install module@version` or plain `go
// build`. The module version comes from Main.Version; commit/date come from
// VCS stamping, which Go only embeds when building from a git checkout (not
// when `go install` fetches the module from the proxy), so they may remain
// unset in that case.
func applyBuildInfo() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	if version == "dev" && info.Main.Version != "" && info.Main.Version != "(devel)" {
		version = info.Main.Version
	}
	var revision, modified string
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			revision = s.Value
		case "vcs.time":
			if date == "unknown" {
				date = s.Value
			}
		case "vcs.modified":
			modified = s.Value
		}
	}
	if commit == "none" && revision != "" {
		if len(revision) > 12 {
			revision = revision[:12]
		}
		if modified == "true" {
			revision += "-dirty"
		}
		commit = revision
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of 17vpn",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("17vpn version %s (commit: %s, built at: %s)\n", version, commit, date)
	},
}

func init() {
	applyBuildInfo()
	rootCmd.Version = version
	rootCmd.AddCommand(versionCmd)
}
