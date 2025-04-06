package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/rinimisini112/tpass/internal/utils"
	"math/rand"
	"strings"
	"time"
)

type PassGenerator struct {
	genLengthInput    textinput.Model
	genIncludeUpper   bool
	genIncludeDigits  bool
	genIncludeSymbols bool
	genResult         string

	focusIndex int
}

const (
	genFocusLength = iota
	genFocusUpper
	genFocusDigits
	genFocusSymbols
)

var lowercase = "abcdefghijklmnopqrstuvwxyz"
var uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
var digits = "0123456789"
var symbols = "!@#$%^&*()_+-=[]{}|;:,.<>?/"

func (g PassGenerator) ParseLength(s string) int {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	if err != nil || n < 1 {
		return 12
	}
	return n
}

func (g PassGenerator) GeneratePassword(length int, includeUpper, includeDigits, includeSymbols bool) string {
	charset := lowercase
	if includeUpper {
		charset += uppercase
	}
	if includeDigits {
		charset += digits
	}
	if includeSymbols {
		charset += symbols
	}

	if charset == "" {
		charset = lowercase
	}

	rand.Seed(time.Now().UnixNano())
	var sb strings.Builder
	for i := 0; i < length; i++ {
		idx := rand.Intn(len(charset))
		sb.WriteByte(charset[idx])
	}
	return sb.String()
}

func (m *Model) RenderPasswordGenerator(underlay string) string {
	g := m.PassGenerator

	var lengthLine string
	if g.focusIndex == genFocusLength {
		lengthLine = generatorFocusStyle.Render("Length: ") + g.genLengthInput.View()
	} else {
		lengthLine = generatorNormalStyle.Render("Length: ") + g.genLengthInput.View()
	}

	upperLine := renderCheckbox("Include Uppercase", g.genIncludeUpper,
		g.focusIndex == genFocusUpper)
	digitsLine := renderCheckbox("Include Digits", g.genIncludeDigits,
		g.focusIndex == genFocusDigits)
	symbolsLine := renderCheckbox("Include Symbols", g.genIncludeSymbols,
		g.focusIndex == genFocusSymbols)

	result := g.genResult
	if result == "" {
		result = "Press [Generate] to create a password!"
	}

	resultLine := lipgloss.NewStyle().Bold(true).Render("Generated Password: ") + result

	content := lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.NewStyle().Bold(true).PaddingBottom(1).Render("PASSWORD GENERATOR"),
		lengthLine,
		upperLine,
		digitsLine,
		symbolsLine,
		lipgloss.NewStyle().PaddingTop(1).Render(resultLine),
		"",
		"Use Up/Down or j/k to move. Left/Right to change length. Space toggles. Enter to activate. ESC to go back.",
	)

	modal := modalStyle.Width(75).Render(content)
	overlayX := (m.Width+4)/2 - modalStyle.GetWidth()/2
	overlayY := m.Height/2 - 8/2

	return utils.PlaceOverlay(overlayX, overlayY, modal, underlay)
}

func renderCheckbox(label string, checked bool, focused bool) string {
	var box string
	if checked {
		box = checkboxOnStyle.Render("[✕]")
	} else {
		box = checkboxOffStyle.Render("[ ]")
	}

	if focused {
		return box + " " + generatorFocusStyle.Render(label)
	} else {
		return box + " " + generatorNormalStyle.Render(label)
	}
}
