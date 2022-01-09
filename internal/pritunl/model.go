package pritunl

type Profile struct {
	ID     string
	Path   string
	Server string
	User   string
}

type Conf struct {
	Name   string `json:"name"`
	Server string `json:"server"`
	User   string `json:"user"`
}

type Connection struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	Timestamp  int64  `json:"timestamp"`
	ServerAddr string `json:"server_addr"`
	ClientAddr string `json:"client_addr"`
}

