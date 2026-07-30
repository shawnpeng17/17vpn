package pritunl

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	unixSocketPath = "/var/run/pritunl.sock"
	daemonTimeout  = 5 * time.Second
	daemonRetries  = 3
	retryBackoff   = 200 * time.Millisecond
)

// daemon talks to the local Pritunl service. Unlike gotunl's makeReq, it
// returns errors instead of calling log.Fatal, so a transient hiccup (e.g. the
// daemon closing an idle socket with EOF) never kills the process.
type daemon struct {
	client  *http.Client
	baseURL string
}

func newDaemon() *daemon {
	if runtime.GOOS == "windows" {
		return &daemon{
			client:  &http.Client{Timeout: daemonTimeout},
			baseURL: "http://localhost:9770/",
		}
	}
	// pritunl uses a unix socket on linux and macOS.
	return &daemon{
		client: &http.Client{
			Timeout: daemonTimeout,
			Transport: &http.Transport{
				// Never reuse a possibly half-closed pooled connection; each
				// request dials the socket fresh.
				DisableKeepAlives: true,
				DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
					var d net.Dialer
					return d.DialContext(ctx, "unix", unixSocketPath)
				},
			},
		},
		baseURL: "http://unix/",
	}
}

// do issues a request to the daemon, retrying a few times on transient errors.
func (d *daemon) do(verb, endpoint, data string) (string, error) {
	var lastErr error
	for attempt := 0; attempt < daemonRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(retryBackoff)
		}
		body, err := d.doOnce(verb, endpoint, data)
		if err == nil {
			return body, nil
		}
		lastErr = err
	}
	return "", lastErr
}

func (d *daemon) doOnce(verb, endpoint, data string) (string, error) {
	req, err := http.NewRequest(verb, d.baseURL+endpoint, strings.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "pritunl")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Auth-Key", authKey())

	res, err := d.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("pritunl daemon returned status %d", res.StatusCode)
	}
	return string(body), nil
}

// authKey reads the daemon auth key fresh each call so a daemon restart that
// rotates the key is picked up automatically. Missing file yields "".
func authKey() string {
	path := "/var/run/pritunl.auth"
	if runtime.GOOS == "windows" {
		path = `c:\ProgramData\Pritunl\auth`
	}
	key, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(key)
}
