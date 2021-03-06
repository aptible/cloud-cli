package common

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// commented unused variables for now
	indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	//subtleIndigo = lipgloss.AdaptiveColor{Light: "#7D79F6", Dark: "#514DC1"}
	//faintRed     = lipgloss.AdaptiveColor{Light: "#FF6F91", Dark: "#C74665"}
	cream       = lipgloss.AdaptiveColor{Light: "#FFFDF5", Dark: "#FFFDF5"}
	yellowGreen = lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#ECFD65"}
	magenta     = lipgloss.Color("#F684FF")
	fuschia     = lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}
	green       = lipgloss.Color("#04B575")
	red         = lipgloss.AdaptiveColor{Light: "#FF4672", Dark: "#ED567A"}
	pink        = lipgloss.AdaptiveColor{Light: "#FF8BA7", Dark: "#893D4E"}
	white       = lipgloss.Color("#FFFFFF")
	lightGrey   = lipgloss.AdaptiveColor{Light: "#BCBCBC", Dark: "#646464"}
	grey        = lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"}
	darkGrey    = lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"}
)

type Styles struct {
	Logo,
	Window,
	Cursor,
	Wrap,
	Paragraph,
	Error,
	Prompt,
	FocusedPrompt,
	SelectionMarker,
	SelectedMenuItem,
	SuccessText,
	InfoText,
	ErrorText,
	Checkmark lipgloss.Style
}

func DefaultStyles() Styles {
	s := Styles{}

	s.Window = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(indigo).
		Padding(1, 2).
		MarginRight(1).
		Width(24)
	s.Logo = lipgloss.NewStyle().
		Foreground(cream).
		Background(indigo).
		Padding(0, 1)
	s.Cursor = lipgloss.NewStyle().Foreground(fuschia)
	s.Wrap = lipgloss.NewStyle().Width(58)
	s.Paragraph = s.Wrap.Copy().Margin(1, 0, 0, 2)
	s.Error = lipgloss.NewStyle().Foreground(red)
	s.Prompt = lipgloss.NewStyle().MarginRight(1).SetString(">")
	s.FocusedPrompt = s.Prompt.Copy().Foreground(fuschia)
	s.SelectionMarker = lipgloss.NewStyle().
		Foreground(fuschia).
		PaddingRight(1).
		SetString(">")
	s.SelectedMenuItem = lipgloss.NewStyle().Foreground(fuschia)
	s.Checkmark = lipgloss.NewStyle().
		SetString("✔").
		Foreground(green)
	s.SuccessText = lipgloss.NewStyle().
		Foreground(green)
	s.InfoText = lipgloss.NewStyle().Foreground(grey)
	s.ErrorText = lipgloss.NewStyle().Foreground(red)

	return s
}

var MainStyles = DefaultStyles()
