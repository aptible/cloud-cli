package loader

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type Model struct {
	Spinner spinner.Model
	Err     error
	Exit    bool
	Text    string
}

func NewModel(text string) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return Model{Spinner: s, Text: text}
}

func (m Model) Init() tea.Cmd {
	return m.Spinner.Tick
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd

	case errMsg:
		m.Err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	}

}

func (m Model) Tick() tea.Msg {
	return m.Spinner.Tick()
}

func (m Model) View() string {
	if m.Err != nil {
		return m.Err.Error()
	}
	str := fmt.Sprintf("%s %s...\n", m.Spinner.View(), m.Text)
	return str
}
