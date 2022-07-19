package fetch

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/aptible/cloud-cli/internal/ui/common"
	"github.com/aptible/cloud-cli/internal/ui/loader"
)

type state int

const (
	submitting state = iota
	success
	quitting
)

type Success struct {
	Result interface{}
}

type errMsg error

type Model struct {
	spinner loader.Model
	Result  interface{}
	Err     error
	io      Fx
	status  state
	styles  common.Styles
}

type Fx func() (dataModel interface{}, error error)

func NewModel(text string, io Fx) Model {
	s := loader.NewModel(text)
	return Model{spinner: s, io: io, status: submitting, styles: common.MainStyles}
}

func create(fx Fx) tea.Cmd {
	return func() tea.Msg {
		res, err := fx()
		if err != nil {
			return err
		}

		return Success{Result: res}
	}
}

func successCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(300 * time.Millisecond)
		return tea.Quit()
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, create(m.io))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.status = quitting
			return m, tea.Quit
		default:
			return m, nil
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case Success:
		m.status = success
		m.Result = msg.Result
		return m, successCmd()

	case errMsg:
		m.Err = msg
		fmt.Printf("Error encountered: %s\n", msg)
		return m, tea.Quit

	default:
		return m, nil
	}

}

func (m Model) View() string {
	if m.Err != nil {
		return m.Err.Error()
	}
	str := ""
	if m.status == submitting {
		str += m.spinner.View()
	} else if m.status == success {
		str += fmt.Sprintf("%s Success!", m.styles.Checkmark.String())
	} else if m.status == quitting {
		return str + "\n"
	}
	return str
}

func Any(model tea.Model) error {
	p := tea.NewProgram(model)
	err := p.Start()
	return err
}

func WithOutput(model tea.Model) (*Model, error) {
	p := tea.NewProgram(model)
	m, err := p.StartReturningModel()
	if err != nil {
		return nil, err
	}

	n := m.(Model)
	if n.Err != nil {
		n.Update(n.Err) // this also quites the program
		os.Exit(1)
	}

	return &n, nil
}
