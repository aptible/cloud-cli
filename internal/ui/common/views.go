package common

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// State is a general UI state used to help style components.
type State int

// UI states.
const (
	StateNormal State = iota
	StateSelected
	StateActive
	StateSpecial
	StateDeleting
)

var lineColors = map[State]lipgloss.TerminalColor{
	StateNormal:   lightGrey,
	StateSelected: magenta,
	StateDeleting: pink,
	StateSpecial:  green,
}

var valStyle = lipgloss.NewStyle().Foreground(indigo)

var (
	helpDivider = lipgloss.NewStyle().
			Foreground(grey).
			Padding(0, 1).
			Render("•")

	helpSection = lipgloss.NewStyle().
			Foreground(darkGrey)
)

// HelpView renders text intended to display at help text, often at the
// bottom of a view.
func HelpView(sections ...string) string {
	var s string
	if len(sections) == 0 {
		return s
	}

	for i := 0; i < len(sections); i++ {
		s += helpSection.Render(sections[i])
		if i < len(sections)-1 {
			s += helpDivider
		}
	}

	return s
}

// VerticalLine return a vertical line colored according to the given state.
func VerticalLine(state State) string {
	return lipgloss.NewStyle().
		SetString("│").
		Foreground(lineColors[state]).
		String()
}

// KeyValueView renders key-value pairs.
func KeyValueView(stuff ...string) string {
	if len(stuff) == 0 {
		return ""
	}

	var (
		s     string
		index int
	)
	for i := 0; i < len(stuff); i++ {
		if i%2 == 0 {
			// even: key
			s += fmt.Sprintf("%s %s: ", VerticalLine(StateNormal), stuff[i])
			continue
		}
		// odd: value
		s += valStyle.Render(stuff[i])
		s += "\n"
		index++
	}

	return strings.TrimSpace(s)
}
