package libasset

import (
	"fmt"
	"log"
	"time"

	cloudapiclient "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/aptible/cloud-cli/config"
	table "github.com/aptible/cloud-cli/lib/op"
	"github.com/aptible/cloud-cli/ui/common"
	"github.com/aptible/cloud-cli/ui/fetch"
	tea "github.com/charmbracelet/bubbletea"
)

type status int

const (
	statusInit status = iota
	statusReady
)

type Model struct {
	config   *config.CloudConfig
	orgId    string
	asset    *cloudapiclient.AssetOutput
	ops      []cloudapiclient.OperationOutput
	fetchOps fetch.Model
	styles   common.Styles
	width    int
	height   int
	status   status
}

func NewDetailModel(config *config.CloudConfig, orgId string, asset *cloudapiclient.AssetOutput) *Model {
	m := &Model{
		orgId:  orgId,
		asset:  asset,
		styles: common.DefaultStyles(),
		status: statusInit,
		config: config,
	}

	return m
}

func RunDetail(config *config.CloudConfig, orgId string, asset *cloudapiclient.AssetOutput) {
	p := tea.NewProgram(NewDetailModel(config, orgId, asset), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.styles.Window.Height(m.height - 4)
		m.styles.Window.Width(m.width - 2)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	case fetch.SuccessMsg:
		m.ops = msg.Result.([]cloudapiclient.OperationOutput)
	}

	switch m.status {
	case statusInit:
		m.status = statusReady
		opMsg := "refreshing"
		m.fetchOps = fetch.NewModelLooper(opMsg, 3*time.Second, func() (interface{}, error) {
			return m.config.Cc.ListOperationsByAsset(m.orgId, m.asset.Id)
		})
		return m, m.fetchOps.Init()
	}

	tmp, cmd := m.fetchOps.Update(message)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	m.fetchOps = tmp.(fetch.Model)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	s := ""
	s += m.styles.Window.Render(m.bioView())
	s += fmt.Sprintf("\n %s", helpView(m))
	return s
}

func helpView(m Model) string {
	var items []string
	items = append(items, "esc: exit")
	return common.HelpView(items...)
}

func (m Model) bioView() string {
	s := m.styles.Logo.Render(GetName(*m.asset))
	s += "\n\n"
	s += common.KeyValueView(
		"Id", m.asset.Id,
		"Asset", m.asset.Asset,
	)
	s += m.opsTableView()
	return s
}

func (m Model) opsTableView() string {
	tbl := table.OpTable(m.ops)
	s := "\n\n\n"
	s += m.styles.Logo.Render("Operations")
	s += "  " + m.fetchOps.View()
	s += "\n"
	s += tbl.View()
	return s
}
