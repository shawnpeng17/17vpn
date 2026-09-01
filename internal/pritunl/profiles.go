package pritunl

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

var serverOrder = map[string]int{
	"DEV":     1,
	"ALPHA":   2,
	"PREPROD": 3,
	"PROD":    4,
}

// profilesDir returns the directory Pritunl stores user profiles in, or ""
// if it can't be determined for this platform.
func profilesDir() string {
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "pritunl", "profiles")
	case "linux":
		return filepath.Join(os.Getenv("HOME"), ".config", "pritunl", "profiles")
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), "pritunl", "profiles")
	}
	return ""
}

// loadProfiles reads every *.conf profile from disk, sorted by environment.
func loadProfiles() []Profile {
	confs, _ := filepath.Glob(filepath.Join(profilesDir(), "*.conf"))

	var profiles []Profile
	for _, path := range confs {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var conf Conf
		_ = json.Unmarshal(data, &conf)
		profiles = append(profiles, Profile{
			ID:     strings.TrimSuffix(filepath.Base(path), ".conf"),
			Path:   path,
			Server: conf.Server,
			User:   conf.User,
		})
	}

	sort.Slice(profiles, func(i, j int) bool {
		return serverOrder[profiles[i].Server] < serverOrder[profiles[j].Server]
	})
	return profiles
}

// profilePayload returns the .ovpn contents for a profile, with the macOS
// keychain secret appended when present. Most profiles have their key
// embedded directly in the .ovpn file and never get a keychain entry — the
// keychain is only populated when Pritunl's client syncs a profile update
// that omits the key, so a missing entry is normal and not an error.
func profilePayload(id string) (string, error) {
	ovpn, err := os.ReadFile(filepath.Join(profilesDir(), id+".ovpn"))
	if err != nil {
		return "", err
	}
	payload := string(ovpn)

	if runtime.GOOS == "darwin" {
		out, err := exec.Command("security", "find-generic-password", "-w", "-s", "pritunl", "-a", id).Output()
		if err == nil {
			if key, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(out))); err == nil {
				payload += "\n" + string(key)
			}
		}
	}
	return payload, nil
}
