package tui

import (
	"errors"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/shawnpeng17/17vpn/internal/pritunl"
)

const (
	idleInterval   = 2 * time.Second
	fastInterval   = 500 * time.Millisecond
	connectTimeout = 30 * time.Second
	failGrace      = 3 * time.Second

	colServer = 12
	colStatus = 18
	colConn   = 12
)

// PasswordFunc returns the password (PIN + OTP) used to connect a profile.
type PasswordFunc func() string

type (
	tickMsg           time.Time
	connectionsMsg    map[string]pritunl.Connection
	connectIssuedMsg  struct{ id, server string }
	disconnectDoneMsg struct{ server string }
)

var (
	headerStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4")) // blue
	selectedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6")) // cyan
	labelStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6")) // cyan
	detailBorder  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("8")).Padding(0, 1)
	helpStyle     = lipgloss.NewStyle().Faint(true)
	infoStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))            // yellow
	successStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2")) // green
	errStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("1")) // red
)

type model struct {
	p        *pritunl.Pritunl
	password PasswordFunc
	profiles []pritunl.Profile
	conns    map[string]pritunl.Connection
	cursor   int
	message  string

	// pending tracks a connection attempt awaiting its result.
	pendingID    string
	pendingSrv   string
	pendingStart time.Time
	deadline     time.Time

	quitting bool
}

// Run launches the interactive TUI. It blocks until the user quits.
func Run(p *pritunl.Pritunl, password PasswordFunc) error {
	profiles := p.Profiles()
	if len(profiles) == 0 {
		return errors.New("No profile found in Pritunl!")
	}

	m := model{
		p:        p,
		password: password,
		profiles: profiles,
		conns:    map[string]pritunl.Connection{},
	}

	_, err := tea.NewProgram(m).Run()
	return err
}

func (m model) Init() tea.Cmd {
	return m.fetchConnections
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.profiles)-1 {
				m.cursor++
			}
		case "enter":
			return m.handleEnter()
		}
	case tickMsg:
		return m, m.fetchConnections
	case connectionsMsg:
		m.conns = msg
		m.evaluatePending()
		return m, tick(m.pollInterval())
	case connectIssuedMsg:
		m.pendingID = msg.id
		m.pendingSrv = msg.server
		m.pendingStart = time.Now()
		m.deadline = m.pendingStart.Add(connectTimeout)
		return m, m.fetchConnections
	case disconnectDoneMsg:
		m.message = infoStyle.Render(fmt.Sprintf("Disconnected %s", msg.server))
		return m, m.fetchConnections
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder
	b.WriteString(m.listView())
	b.WriteString("\n")
	b.WriteString(m.detailView())
	b.WriteString("\n\n")
	if m.message != "" {
		b.WriteString(m.message)
		b.WriteString("\n")
	}
	b.WriteString(helpStyle.Render("↑/↓ or j/k move · enter connect/disconnect · q quit"))
	b.WriteString("\n")
	return b.String()
}

func (m model) listView() string {
	var b strings.Builder
	b.WriteString("  ")
	b.WriteString(headerStyle.Render(cell("SERVER", colServer) + cell("STATUS", colStatus) + cell("CONNECTED", colConn)))
	b.WriteString("\n")

	for i, prof := range m.profiles {
		status := "disconnected"
		connected := ""
		if c, ok := m.conns[prof.ID]; ok {
			if c.Status != "" {
				status = c.Status
			}
			if c.Timestamp > 0 {
				connected = formatDuration(time.Since(time.Unix(c.Timestamp, 0)))
			}
		}

		pointer := "  "
		server := prof.Server
		if i == m.cursor {
			pointer = selectedStyle.Render("▶ ")
			server = selectedStyle.Render(server)
		}

		b.WriteString(pointer)
		b.WriteString(cell(server, colServer))
		b.WriteString(cell(statusBadge(status), colStatus))
		b.WriteString(cell(connected, colConn))
		b.WriteString("\n")
	}
	return b.String()
}

// cell pads (or truncates) a possibly-styled string to a fixed visible width.
func cell(s string, w int) string {
	return lipgloss.NewStyle().Width(w).Render(s)
}

// handleEnter toggles the selected profile: disconnect it if connected,
// otherwise disconnect every other connection and connect it.
func (m model) handleEnter() (tea.Model, tea.Cmd) {
	if m.cursor < 0 || m.cursor >= len(m.profiles) {
		return m, nil
	}
	target := m.profiles[m.cursor]

	if m.conns[target.ID].Status == "connected" {
		m.message = infoStyle.Render(fmt.Sprintf("Disconnecting %s...", target.Server))
		return m, m.disconnectCmd(target)
	}

	m.message = infoStyle.Render(fmt.Sprintf("Connecting %s...", target.Server))
	return m, m.connectCmd(target)
}

func (m model) fetchConnections() tea.Msg {
	return connectionsMsg(m.p.Connections())
}

func tick(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg { return tickMsg(t) })
}

// connectCmd disconnects all other active connections (mirroring the CLI's
// single-VPN behaviour), then connects the target. It runs off the UI loop.
func (m model) connectCmd(target pritunl.Profile) tea.Cmd {
	p := m.p
	pass := m.password
	profiles := m.profiles
	conns := m.conns
	return func() tea.Msg {
		for _, prof := range profiles {
			if prof.ID == target.ID {
				continue
			}
			if _, ok := conns[prof.ID]; ok {
				p.Disconnect(prof.ID)
				time.Sleep(time.Second)
			}
		}
		p.Connect(target.ID, pass())
		return connectIssuedMsg{id: target.ID, server: target.Server}
	}
}

func (m model) disconnectCmd(target pritunl.Profile) tea.Cmd {
	p := m.p
	return func() tea.Msg {
		p.Disconnect(target.ID)
		return disconnectDoneMsg{server: target.Server}
	}
}

// evaluatePending resolves an in-flight connection attempt into a success,
// failure, or timeout message based on the freshly fetched connection state.
func (m *model) evaluatePending() {
	if m.pendingID == "" {
		return
	}

	status := strings.ToLower(m.conns[m.pendingID].Status)
	switch {
	case status == "connected":
		m.message = successStyle.Render(fmt.Sprintf("Connected %s!", m.pendingSrv))
		m.pendingID = ""
	case status == "" && time.Since(m.pendingStart) > failGrace:
		m.message = errStyle.Render(fmt.Sprintf("Connect %s failed!", m.pendingSrv))
		m.pendingID = ""
	case time.Now().After(m.deadline):
		m.message = errStyle.Render(fmt.Sprintf("Connect %s timeout!", m.pendingSrv))
		m.pendingID = ""
	}
}

// pollInterval speeds up while a connection is transitioning and relaxes
// back to the idle cadence once everything is settled.
func (m model) pollInterval() time.Duration {
	if m.pendingID != "" {
		return fastInterval
	}
	for _, c := range m.conns {
		switch strings.ToLower(c.Status) {
		case "connecting", "disconnecting":
			return fastInterval
		}
	}
	return idleInterval
}

func (m model) detailView() string {
	if m.cursor < 0 || m.cursor >= len(m.profiles) {
		return ""
	}
	prof := m.profiles[m.cursor]

	client, server := "-", "-"
	if c, ok := m.conns[prof.ID]; ok {
		if c.ClientAddr != "" {
			client = c.ClientAddr
		}
		if c.ServerAddr != "" {
			server = c.ServerAddr
		}
	}

	row := func(label, val string) string {
		return labelStyle.Render(fmt.Sprintf("%-10s", label)) + val
	}
	content := strings.Join([]string{
		row("User", prof.User),
		row("Client IP", client),
		row("Server IP", server),
	}, "\n")
	return detailBorder.Render(content)
}

func statusBadge(status string) string {
	upper := strings.ToUpper(status)
	switch upper {
	case "CONNECTED":
		return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2")).Render(upper) // green
	case "CONNECTING":
		return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("3")).Render(upper + "...") // yellow
	case "DISCONNECTING":
		return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("8")).Render(upper + "...") // gray
	case "DISCONNECTED":
		return lipgloss.NewStyle().Faint(true).Render(upper)
	default:
		return lipgloss.NewStyle().Faint(true).Render("UNKNOWN")
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
