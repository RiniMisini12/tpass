package tui

import "github.com/charmbracelet/lipgloss"

var (
	appStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#8839ef")).
			Margin(1, 2).
			Padding(1, 2)

	leftPaneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#8839ef")).
			MarginRight(1).
			Padding(1, 2)

	rightPaneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#8839ef")).
			MarginLeft(1).
			Padding(1, 2)

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8839ef")).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			PaddingLeft(2)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff0000")).
			PaddingLeft(2)

	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#8839ef")).
			Padding(0, 1)

	modalStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#8839ef")).
			Padding(1, 2).
			Width(60)

	generatorFocusStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#8839ef"))

	generatorNormalStyle = lipgloss.NewStyle()

	checkboxOnStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#8839ef"))

	checkboxOffStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff"))
)
