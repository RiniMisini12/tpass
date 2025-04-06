package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rinimisini112/tpass/internal/store"
	"os"
	"time"
)

func DeleteActions(msg tea.Msg, m *Model) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch keyMsg.String() {
	case keyEnter:
		if m.DeleteIndex >= 0 && m.DeleteIndex < len(m.Passwords) {
			m.Passwords = append(m.Passwords[:m.DeleteIndex], m.Passwords[m.DeleteIndex+1:]...)
			m.Store.Entries = m.Passwords
			err := store.SaveStore(m.Store, store.GetMasterPassword())
			if err != nil {
				_, err := fmt.Fprintf(os.Stderr, "Error saving store: %v\n", err)
				if err != nil {
					return nil, nil
				}
				os.Exit(1)
			}
			if m.DeleteIndex >= len(m.Passwords) {
				m.SelectedIndex = len(m.Passwords) - 1
			} else {
				m.SelectedIndex = m.DeleteIndex
			}
			m.State = StateList
			m.Status = "Password deleted!"
			return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
				return clearStatusMsg{}
			})
		}
		m.State = StateList
		return m, nil

	case keyQuit:
		m.State = StateList
		return m, nil
	}

	return m, nil
}
