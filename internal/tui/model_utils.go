package tui

import "github.com/atotto/clipboard"

func (m *Model) CopyPassword(password string) {
	if err := clipboard.WriteAll(password); err != nil {
		panic(err)
	}
}

func (m *Model) FocusField(index int) {
	m.FocusedField = index
	m.NameInput.Blur()
	m.PassInput.Blur()
	m.SearchInput.Blur()

	switch m.FocusedField {
	case 0:
		m.NameInput.Focus()
	case 1:
		m.PassInput.Focus()
	case 2:
		m.SearchInput.Focus()
	}
}
