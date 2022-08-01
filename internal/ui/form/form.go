package form

import (
	"fmt"

	"github.com/aptible/cloud-cli/internal/common"
	uiCommon "github.com/aptible/cloud-cli/internal/ui/common"
	"github.com/aptible/cloud-cli/internal/ui/loader"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

type loadedOptionsMsg struct {
	Options []list.Item
}

type valueEnteredMsg struct {
	Value string
}

type status int

const (
	statusInit status = iota
	statusReady
	statusLoadingOptions
	statusValueEntered
	statusUserInput
)

type Model struct {
	styles   uiCommon.Styles
	config   *common.CloudConfig
	schema   *SubSchema
	list     list.Model
	input    textinput.Model
	spinner  loader.Model
	status   status
	metaDesc string
	Result   string
	Err      error
}

func NewModel(config *common.CloudConfig, schema *SubSchema) *Model {
	ti := textinput.New()
	ti.CharLimit = 156
	ti.Width = 50

	model := &Model{
		styles:  uiCommon.DefaultStyles(),
		spinner: loader.NewModel("fetching resources"),
		config:  config,
		schema:  schema,
		status:  statusInit,
		list:    list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		input:   ti,
	}
	model.list.Title = schema.Title
	return model
}

func (m Model) fetchOptions() tea.Cmd {
	return func() tea.Msg {
		options, err := m.schema.LoadOptions(m.config)
		if err != nil {
			return err
		}

		return loadedOptionsMsg{Options: options}
	}
}

func valueEntered(val string) tea.Cmd {
	return func() tea.Msg {
		return valueEnteredMsg{Value: val}
	}
}

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.schema.Type == "input" {
				return m, valueEntered(m.input.Value())
			} else if m.schema.Type == "select" {
				// we ask users to press enter when applying a filter
				// which means we need to check to make sure that's not
				// why they hit enter
				if !m.list.SettingFilter() {
					val := ""
					if i, ok := m.list.SelectedItem().(FormOption); ok {
						val = i.Value
					}
					return m, valueEntered(val)
				}
			}
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width-5, 20)
	case loadedOptionsMsg:
		if len(msg.Options) == 1 {
			val := ""
			if i, ok := msg.Options[0].(FormOption); ok {
				val = i.Value
			}
			m.metaDesc = " (only option available)"
			return m, valueEntered(val)
		} else {
			m.status = statusReady
			m.list.SetItems(msg.Options)
		}
	case valueEnteredMsg:
		m.status = statusValueEntered
		m.Result = msg.Value
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case errMsg:
		m.Err = msg
		return m, nil
	}

	switch m.status {
	case statusInit:
		m.status = statusReady
		if m.schema.Type == "select" && m.schema.LoadOptions != nil {
			m.status = statusLoadingOptions
			return m, m.fetchOptions()
		}
		if m.schema.Type == "input" {
			m.input.Focus()
			m.status = statusUserInput
		}
	}

	switch m.schema.Type {
	case "input":
		m.input, cmd = m.input.Update(message)
	case "select":
		m.list, cmd = m.list.Update(message)
	}
	return m, cmd
}

func (m Model) View() string {
	if m.Err != nil {
		return m.Err.Error()
	}

	s := ""

	if m.status == statusReady && m.schema.Type == "select" {
		s += fmt.Sprintf("\n%s", m.list.View())
	} else if m.status == statusLoadingOptions {
		s += m.spinner.View()
	} else if m.status == statusValueEntered {
		s += fmt.Sprintf(
			"%s: %s%s\n",
			m.schema.Title,
			m.styles.SuccessText.Render(m.Result),
			m.styles.InfoText.Render(m.metaDesc),
		)
	} else if m.status == statusUserInput {
		s += fmt.Sprintf("%s: %s", m.schema.Title, m.input.View())
	}

	return s
}

func Run(model *Model) (string, error) {
	p := tea.NewProgram(model)
	m, err := p.StartReturningModel()
	if err != nil {
		return "", err
	}

	switch n := m.(type) {
	case Model:
		if n.Err != nil {
			n.Update(n.Err) // this also quits the program
			return "", n.Err
		}

		return n.Result, nil
	default:
		return "", fmt.Errorf("woops")
	}
}
