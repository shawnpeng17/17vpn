package pritunl

import (
	"encoding/json"
	"sort"

	"github.com/cghdev/gotunl/pkg/gotunl"
)

var serverOrder = map[string]int{
	"DEV":   1,
	"ALPHA": 2,
	"STA":   3,
	"PROD":  4,
}

type Pritunl struct {
	gotunl *gotunl.Gotunl
}

func New() *Pritunl {
	return &Pritunl{
		gotunl: gotunl.New(),
	}
}

func (p *Pritunl) Profiles() []Profile {
	var profiles []Profile
	for id, profile := range p.gotunl.Profiles {
		var conf Conf
		_ = json.Unmarshal([]byte(profile.Conf), &conf)
		profiles = append(profiles, Profile{
			ID:     id,
			Path:   profile.Path,
			Server: conf.Server,
			User:   conf.User,
		})
	}
	sort.Slice(profiles, func(i, j int) bool {
		return serverOrder[profiles[i].Server] < serverOrder[profiles[j].Server]
	})
	return profiles
}

func (p *Pritunl) Connections() map[string]Connection {
	var conns map[string]Connection
	connStr := p.gotunl.GetConnections()
	_ = json.Unmarshal([]byte(connStr), &conns)
	return conns
}

func (p *Pritunl) Connect(id, password string) {
	p.gotunl.ConnectProfile(id, "pritunl", password)
}

func (p *Pritunl) Disconnect(id string) {
	p.gotunl.DisconnectProfile(id)
}

func (p *Pritunl) DisconnectAll() {
	p.gotunl.StopConnections()
}
