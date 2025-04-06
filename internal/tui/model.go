package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rinimisini112/tpass/internal/store"
)

type state int

const (
	StateList state = iota
	StateFilter
	StateAdd
	StateEdit
	StateConfirmDelete
	StateGeneratePassword
)

type clearStatusMsg struct{}

type clearErrorMsg struct{}

type Model struct {
	State         state
	Store         store.Store
	PassGenerator PassGenerator

	Width  int
	Height int

	SearchInput   textinput.Model
	Passwords     []store.PasswordEntry
	SelectedIndex int

	FocusedField int
	NameInput    textinput.Model
	PassInput    textinput.Model

	Status      string
	Error       string
	ShowInfo    bool
	ShowPreview bool

	DeleteIndex int
}

func InitialModel(store store.Store) *Model {
	search := textinput.New()
	search.Placeholder = "(Press / to Search)"
	search.CharLimit = 200

	name := textinput.New()
	name.Placeholder = "Name"
	name.CharLimit = 200

	pass := textinput.New()
	pass.Placeholder = "Password"
	pass.CharLimit = 200

	m := &Model{
		State:         StateList,
		Store:         store,
		Passwords:     store.Entries,
		SelectedIndex: 0,
		SearchInput:   search,
		NameInput:     name,
		PassInput:     pass,
		Status:        "",
		Error:         "",
		ShowPreview:   false,
		PassGenerator: PassGenerator{},
	}
	m.PassGenerator.focusIndex = genFocusLength
	m.PassGenerator.genLengthInput = textinput.New()
	m.PassGenerator.genLengthInput.Placeholder = "Length (e.g. 12)"
	m.PassGenerator.genLengthInput.SetValue("12")

	m.PassGenerator.genIncludeUpper = true
	m.PassGenerator.genIncludeDigits = true
	m.PassGenerator.genIncludeSymbols = false
	m.PassGenerator.genResult = ""
	m.FocusField(0)
	return m
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width - 10
		m.Height = msg.Height - 8
		return m, nil

	case clearStatusMsg:
		m.Status = ""
		return m, nil

	case clearErrorMsg:
		m.Error = ""
		return m, nil

	case tea.KeyMsg:
		if m.ShowInfo {
			if msg.String() == keyQuit {
				m.ShowInfo = false
			}
			return m, nil
		}
		switch m.State {
		case StateList, StateFilter:
			return ListActions(msg, m)

		case StateAdd, StateEdit:
			AddEditActions(msg, m)

		case StateConfirmDelete:
			DeleteActions(msg, m)

		case StateGeneratePassword:
			PassGeneratorActions(msg, m)
		}
	}

	if m.State == StateAdd {
		switch m.FocusedField {
		case 0:
			m.NameInput, cmd = m.NameInput.Update(msg)
		case 1:
			m.PassInput, cmd = m.PassInput.Update(msg)
		}
	}

	if m.State == StateEdit {
		switch m.FocusedField {
		case 0:
			m.NameInput, cmd = m.NameInput.Update(msg)
		case 1:
			m.PassInput, cmd = m.PassInput.Update(msg)
		}
	}

	if m.State == StateFilter {
		m.SearchInput, cmd = m.SearchInput.Update(msg)
		m.SelectedIndex = 0
	}

	return m, cmd
}

func (m *Model) View() string {
	if m.Width <= 0 || m.Height <= 0 {
		return "Loading..."
	}

	header := m.RenderHeader()
	outerWidth := m.Width - appStyle.GetPaddingLeft() - appStyle.GetPaddingRight() - 2
	outerHeight := m.Height - appStyle.GetPaddingTop() - appStyle.GetPaddingBottom() - 2

	leftWidth := outerWidth / 2
	rightWidth := outerWidth - leftWidth

	var leftView, rightView string

	switch m.State {
	case StateList, StateFilter, StateConfirmDelete, StateGeneratePassword:
		leftView = m.RenderList(leftWidth)
		rightView = m.RenderPreview(m.ShowPreview)
	case StateAdd:
		leftView = m.RenderList(leftWidth)
		rightView = m.RenderAddForm(rightWidth, false)
	case StateEdit:
		leftView = m.RenderList(leftWidth)
		rightView = m.RenderAddForm(rightWidth, true)
	}

	leftPane := leftPaneStyle.Width(leftWidth).Height(outerHeight).Render(leftView)
	rightPane := rightPaneStyle.Width(rightWidth).Height(outerHeight).Render(rightView)
	body := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)

	ui := lipgloss.JoinVertical(lipgloss.Top, header, body)
	finalUI := appStyle.Render(ui)

	if m.State == StateGeneratePassword {
		return m.RenderPasswordGenerator(finalUI)
	}

	if m.State == StateConfirmDelete {
		return m.RenderConfirmDeleteModal(finalUI)
	}

	if m.ShowInfo {
		return m.RenderHelpMenu(finalUI)
	}

	return finalUI
}
