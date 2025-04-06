package tui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/rinimisini112/tpass/internal/store"
	"github.com/rinimisini112/tpass/internal/utils"
	"strings"
)

func (m *Model) FilteredPasswords() []store.PasswordEntry {
	query := strings.ToLower(m.SearchInput.Value())
	if query == "" {
		return m.Passwords
	}
	var filtered []store.PasswordEntry
	for _, p := range m.Passwords {
		if strings.Contains(strings.ToLower(p.Name), query) {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

func (m *Model) RenderList(width int) string {
	filtered := m.FilteredPasswords()
	if m.SelectedIndex >= len(filtered) {
		m.SelectedIndex = 0
	}
	s := "Passwords (Press Up/Down to select, 🅰 to add, CTRL + 🅲 to quit)\n"
	s += inputStyle.Width(width-6).Render(m.SearchInput.View()) + "\n\n"
	if len(filtered) > 0 {
		for i, p := range filtered {
			prefix := "  "
			if i == m.SelectedIndex {
				prefix = "> "
			}
			s += fmt.Sprintf("%s%s\n", prefix, p.Name)
		}
	} else {
		s += "All who wander are not lost but what ur looking for is nowhere to be found.\n"
	}
	s += "\nPress 🅲 to copy (yank) selected password."
	return s
}

func (m *Model) RenderPreview(show bool) string {
	filtered := m.FilteredPasswords()
	if len(filtered) == 0 {
		return "No passwords match your search.\n"
	}
	if m.SelectedIndex >= len(filtered) {
		return "No selection.\n"
	}

	p := filtered[m.SelectedIndex]
	password := p.Password
	if !show {
		password = strings.Repeat("*", len(password))
	}

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8839ef")).
		Bold(true).
		PaddingBottom(1)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#c0a0eb")).
		PaddingBottom(1).
		Bold(true)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff"))

	headerText := "Password Preview"
	if show {
		headerText += " - Showing password"
	}

	header := headerStyle.Render(headerText)

	nameLine := lipgloss.JoinHorizontal(lipgloss.Top,
		labelStyle.Render("Name: "),
		valueStyle.Render(p.Name),
	)

	passwordLine := lipgloss.JoinHorizontal(lipgloss.Top,
		labelStyle.Render("Password: "),
		valueStyle.Render(password),
	)

	return lipgloss.JoinVertical(lipgloss.Left, header, nameLine, passwordLine)
}

func (m *Model) RenderAddForm(width int, isEdit bool) string {
	s := "Add a new password:\n\n"
	if isEdit {
		s = "Password edit:\n\n"
	}
	s += "Name:\n" + inputStyle.Width(width-6).Render(m.NameInput.View()) + "\n\n"
	s += "Password:\n" + inputStyle.Width(width-6).Render(m.PassInput.View()) + "\n\n"
	s += "(Press Enter to save, Ctrl+🅲 to cancel, Tab/Up/Down to switch fields)"
	return s
}

func (m *Model) RenderConfirmDeleteModal(finalUI string) string {
	confirmText := "Are you sure you want to delete this entry?\n\n" +
		"Press Enter to confirm, Ctrl+c to cancel."
	modalView := lipgloss.NewStyle().
		Padding(2, 4).
		Width(55).
		Height(8).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FF5F5F")).
		Render(confirmText)

	overlayX := (m.Width+4)/2 - 55/2
	overlayY := m.Height/2 - 8/2

	return utils.PlaceOverlay(overlayX, overlayY, modalView, finalUI)
}

func (m *Model) RenderHelpMenu(finalUI string) string {
	modalView := lipgloss.NewStyle().
		Padding(2, 4).
		Width(55).
		Height(8).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Render("This is a modal popup!")

	overlayX := (m.Width+4)/2 - 55/2
	overlayY := m.Height/2 - 8/2
	return utils.PlaceOverlay(overlayX, overlayY, modalView, finalUI)
}

func (m *Model) RenderHeader() string {
	navItems := "🅼  For Password Generator"
	navStyled := lipgloss.NewStyle().
		Bold(true).
		PaddingBottom(1).
		Render(navItems)

	var leftHalf string
	header := headerStyle.Render("TPASS Terminal Password Manager")

	if m.Error != "" {
		leftHalf = lipgloss.JoinHorizontal(lipgloss.Left, header, "  ", errorStyle.Render(m.Error))
	} else if m.Status != "" {
		leftHalf = lipgloss.JoinHorizontal(lipgloss.Left, header, "  ", statusStyle.Render(m.Status))
	} else {
		leftHalf = header
	}

	leftColumnWidth := (m.Width * 2) / 3
	rightColumnWidth := m.Width - leftColumnWidth

	leftCol := lipgloss.NewStyle().
		Width(leftColumnWidth).
		Align(lipgloss.Left).
		Render(leftHalf)

	rightCol := lipgloss.NewStyle().
		Width(rightColumnWidth).
		Align(lipgloss.Right).
		Render(navStyled)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftCol, rightCol)
}
