package fetch

import (
	"fmt"
	"os"
	"time"

	"github.com/aptible/cloud-cli/ui/common"
	"github.com/aptible/cloud-cli/ui/loader"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	ready state = iota
	submitting
	success
	quitting
)

type SuccessMsg struct {
	Result interface{}
}

type fetchMsg struct{}

type errMsg error

type Model struct {
	spinner loader.Model
	Result  interface{}
	Err     error
	io      Fx
	status  state
	styles  common.Styles
	Loop    time.Duration
}

type Fx func() (dataModel interface{}, error error)

func NewModel(text string, io Fx) Model {
	s := loader.NewModel(text)
	return Model{
		spinner: s,
		io:      io,
		status:  ready,
		styles:  common.MainStyles,
	}
}

func NewModelLooper(text string, loop time.Duration, io Fx) Model {
	mdl := NewModel(text, io)
	mdl.Loop = loop
	return mdl
}

func create(fx Fx) tea.Cmd {
	return func() tea.Msg {
		res, err := fx()
		if err != nil {
			return err
		}

		return SuccessMsg{Result: res}
	}
}

func (m Model) FetchCmd() tea.Cmd {
	return func() tea.Msg {
		return fetchMsg{}
	}
}

func (m Model) successCmd() tea.Cmd {
	return func() tea.Msg {
		if m.Loop == 0 {
			time.Sleep(300 * time.Millisecond)
			return tea.Quit()
		}

		time.Sleep(m.Loop)
		return fetchMsg{}
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.FetchCmd())
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

	case fetchMsg:
		m.status = submitting
		return m, create(m.io)

	case SuccessMsg:
		m.status = success
		m.Result = msg.Result
		return m, m.successCmd()

	case errMsg:
		m.Err = msg
		fmt.Printf("Error encountered: %s\n", msg)
		return m, tea.Quit
	}

	return m, nil
}

func (m Model) View() string {
	if m.Err != nil {
		return m.Err.Error()
	}
	str := ""
	if m.status == submitting {
		str += m.spinner.View()
	} else if m.status == success {
		str += fmt.Sprintf("%s success!", m.styles.Checkmark.String())
	} else if m.status == quitting {
		str += "\n"
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
