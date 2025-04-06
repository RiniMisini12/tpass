package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

func PassGeneratorActions(msg tea.Msg, m *Model) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch keyMsg.String() {
	case keyUp:
		fallthrough
	case keyUpAlt:
		if m.PassGenerator.focusIndex > 0 {
			m.PassGenerator.focusIndex--
		}
		return m, nil
	case keyEnter:
		length := m.PassGenerator.ParseLength(m.PassGenerator.genLengthInput.Value())
		m.PassGenerator.genResult = m.PassGenerator.GeneratePassword(
			length,
			m.PassGenerator.genIncludeUpper,
			m.PassGenerator.genIncludeDigits,
			m.PassGenerator.genIncludeSymbols,
		)

	case keyCopy:
		m.CopyPassword(m.PassGenerator.genResult)
		m.Status = "Generated password copied to clipboard!"
		return m, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
			return clearStatusMsg{}
		})

	case keyAdd:
		m.NameInput.SetValue("")
		m.PassInput.SetValue(m.PassGenerator.genResult)
		m.State = StateAdd
		return m, nil
	case keyDown:
		fallthrough
	case keyDownAlt:
		if m.PassGenerator.focusIndex < genFocusSymbols {
			m.PassGenerator.focusIndex++
		}
		return m, nil
	case keyLeft:
		if m.PassGenerator.focusIndex == genFocusLength {
			length := m.PassGenerator.ParseLength(m.PassGenerator.genLengthInput.Value())
			if length > 1 {
				length--
				m.PassGenerator.genLengthInput.SetValue(fmt.Sprintf("%d", length))
			}
		}
		return m, nil

	case keyRight:
		if m.PassGenerator.focusIndex == genFocusLength {
			length := m.PassGenerator.ParseLength(m.PassGenerator.genLengthInput.Value())
			length++
			m.PassGenerator.genLengthInput.SetValue(fmt.Sprintf("%d", length))
		}
		return m, nil

	case keySpace:
		switch m.PassGenerator.focusIndex {
		case genFocusUpper:
			m.PassGenerator.genIncludeUpper = !m.PassGenerator.genIncludeUpper
		case genFocusDigits:
			m.PassGenerator.genIncludeDigits = !m.PassGenerator.genIncludeDigits
		case genFocusSymbols:
			m.PassGenerator.genIncludeSymbols = !m.PassGenerator.genIncludeSymbols
		}
		return m, nil

	case keyQuit:
		m.State = StateList
		return m, nil
	}

	if m.PassGenerator.focusIndex == genFocusLength {
		var tiCmd tea.Cmd
		m.PassGenerator.genLengthInput, tiCmd = m.PassGenerator.genLengthInput.Update(msg)
		return m, tiCmd
	}

	return m, nil
}
