package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rinimisini112/tpass/internal/store"
	"os"
	"time"
)

func AddEditActions(msg tea.Msg, m *Model) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch m.State {
	case StateAdd:
		switch keyMsg.String() {
		case keyQuit:
			m.State = StateList
			return m, nil
		case keyEnter:
			if m.NameInput.Value() == "" || m.PassInput.Value() == "" {
				m.Error = "Inputs cannot be empty!"
				return m, tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
					return clearErrorMsg{}
				})
			}

			newPass := store.PasswordEntry{
				Name:     m.NameInput.Value(),
				Password: m.PassInput.Value(),
			}

			m.Store.Entries = append(m.Store.Entries, newPass)
			if err := store.SaveStore(m.Store, store.GetMasterPassword()); err != nil {
				_, err := fmt.Fprintf(os.Stderr, "Error saving store: %v\n", err)
				if err != nil {
					return nil, nil
				}
				os.Exit(1)
			}

			m.Passwords = m.Store.Entries
			m.NameInput.SetValue("")
			m.PassInput.SetValue("")
			m.SelectedIndex = len(m.Passwords) - 1
			m.State = StateList
			m.Status = "Password added successfully!"
			return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
				return clearStatusMsg{}
			})
			return m, nil
		case keyTab, keyDown:
			m.FocusField((m.FocusedField + 1) % 2)
			return m, nil
		case keyUp:
			m.FocusField((m.FocusedField + 1) % 2)
			return m, nil
		}
	case StateEdit:
		switch keyMsg.String() {
		case keyQuit:
			m.State = StateList
			return m, nil
		case keyEnter:
			if m.NameInput.Value() == "" || m.PassInput.Value() == "" {
				m.Error = "Inputs cannot be empty!"
				return m, tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
					return clearErrorMsg{}
				})
			}
			m.Passwords[m.SelectedIndex] = store.PasswordEntry{
				Name:     m.NameInput.Value(),
				Password: m.PassInput.Value(),
			}
			m.Store.Entries = m.Passwords
			if err := store.SaveStore(m.Store, store.GetMasterPassword()); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error saving store: %v\n", err)
				os.Exit(1)
			}
			m.State = StateList
			m.Status = "Password updated successfully!"
			return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
				return clearStatusMsg{}
			})
		case keyTab, keyDown:
			m.FocusField((m.FocusedField + 1) % 2)
			return m, nil
		case keyUp:
			m.FocusField((m.FocusedField + 1) % 2)
			return m, nil
		}
	}

	return m, nil
}
