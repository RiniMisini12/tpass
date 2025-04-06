package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

func ListActions(msg tea.Msg, m *Model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch m.State {
	case StateList:
		switch keyMsg.String() {
		case keyHelp:
			m.ShowInfo = !m.ShowInfo
			return m, nil
		case keyQuit:
			return m, tea.Quit

		case keyAdd:
			m.NameInput.SetValue("")
			m.PassInput.SetValue("")
			m.State = StateAdd
			m.FocusField(0)
			return m, nil
		case keyPassGenerate:
			m.State = StateGeneratePassword
			return m, nil
		case keyShowPassword:
			m.ShowPreview = !m.ShowPreview
			return m, nil

		case keyEdit:
			if len(m.Passwords) > 0 {
				entry := m.Passwords[m.SelectedIndex]
				m.NameInput.SetValue(entry.Name)
				m.PassInput.SetValue(entry.Password)
				m.State = StateEdit
				m.FocusField(0)
			}
			return m, nil

		case keyToggleSearch:
			if m.SearchInput.Focused() {
				m.SearchInput.Blur()
				m.State = StateList
			} else {
				m.FocusField(2)
				m.State = StateFilter
			}
			return m, nil

		case keyUp:
			if m.SelectedIndex > 0 {
				m.SelectedIndex--
			}
			return m, nil
		case keyDown:
			if m.SelectedIndex < len(m.FilteredPasswords())-1 {
				m.SelectedIndex++
			}
			return m, nil
		case keyCopy:
			if len(m.Passwords) == 0 || m.SelectedIndex >= len(m.Passwords) {
				return m, nil
			}
			filtered := m.FilteredPasswords()
			if m.SelectedIndex >= len(filtered) {
				return m, nil
			}
			p := filtered[m.SelectedIndex]

			m.CopyPassword(p.Password)
			m.Status = "Password copied to clipboard!"

			return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
				return clearStatusMsg{}
			})
		case keyCopyAsString:
			if len(m.Passwords) == 0 || m.SelectedIndex >= len(m.Passwords) {
				return m, nil
			}
			filtered := m.FilteredPasswords()
			if m.SelectedIndex >= len(filtered) {
				return m, nil
			}
			p := filtered[m.SelectedIndex]
			password := "\"" + p.Password + "\""
			m.CopyPassword(password)
			m.Status = "Password copied to clipboard!"
			return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
				return clearStatusMsg{}
			})
		case keyDelete:
			if len(m.Passwords) > 0 {
				m.DeleteIndex = m.SelectedIndex
				m.State = StateConfirmDelete
			}
			return m, nil
		}
	case StateFilter:
		switch keyMsg.String() {
		case keyQuit, keyEnter:
			m.SearchInput.Blur()
			m.State = StateList
			return m, nil
		}

		m.SearchInput, cmd = m.SearchInput.Update(msg)
		m.SelectedIndex = 0
		return m, cmd
	}

	return m, nil
}
