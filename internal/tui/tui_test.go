package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shawnpeng17/17vpn/internal/pritunl"
)

func TestCursorCycle(t *testing.T) {
	m := model{
		profiles: []pritunl.Profile{
			{ID: "1", Server: "Server 1"},
			{ID: "2", Server: "Server 2"},
			{ID: "3", Server: "Server 3"},
		},
		cursor: 0,
	}

	// Press "up" from index 0 -> should cycle to index 2 (last item)
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updatedModel.(model)
	if m.cursor != 2 {
		t.Errorf("expected cursor 2 after pressing up from 0, got %d", m.cursor)
	}

	// Press "down" from index 2 -> should cycle to index 0 (first item)
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updatedModel.(model)
	if m.cursor != 0 {
		t.Errorf("expected cursor 0 after pressing down from 2, got %d", m.cursor)
	}

	// Press "k" (vim up) from index 0 -> index 2
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updatedModel.(model)
	if m.cursor != 2 {
		t.Errorf("expected cursor 2 after pressing 'k' from 0, got %d", m.cursor)
	}

	// Press "j" (vim down) from index 2 -> index 0
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updatedModel.(model)
	if m.cursor != 0 {
		t.Errorf("expected cursor 0 after pressing 'j' from 2, got %d", m.cursor)
	}
}
