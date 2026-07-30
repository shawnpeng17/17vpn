package pritunl

import (
	"encoding/json"
)

type Pritunl struct {
	daemon *daemon
}

func New() *Pritunl {
	return &Pritunl{
		daemon: newDaemon(),
	}
}

func (p *Pritunl) Profiles() []Profile {
	return loadProfiles()
}

func (p *Pritunl) Connections() (map[string]Connection, error) {
	body, err := p.daemon.do("GET", "profile", "")
	if err != nil {
		return nil, err
	}
	var conns map[string]Connection
	if err := json.Unmarshal([]byte(body), &conns); err != nil {
		return nil, err
	}
	return conns, nil
}

func (p *Pritunl) Connect(id, password string) error {
	ovpn, err := profilePayload(id)
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"id":        id,
		"reconnect": true,
		"timeout":   true,
		"username":  "pritunl",
		"password":  password,
		"data":      ovpn,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = p.daemon.do("POST", "profile", string(data))
	return err
}

func (p *Pritunl) Disconnect(id string) error {
	data, err := json.Marshal(map[string]string{"id": id})
	if err != nil {
		return err
	}
	_, err = p.daemon.do("DELETE", "profile", string(data))
	return err
}

func (p *Pritunl) DisconnectAll() error {
	_, err := p.daemon.do("POST", "stop", "")
	return err
}
