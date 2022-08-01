package asset

import (
	"fmt"
	"log"
	"time"

	cloudapiclient "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/aptible/cloud-cli/internal/common"
	uiCommon "github.com/aptible/cloud-cli/internal/ui/common"
	"github.com/aptible/cloud-cli/internal/ui/fetch"
	render "github.com/aptible/cloud-cli/table"
	tea "github.com/charmbracelet/bubbletea"
)

type status int

const (
	statusInit status = iota
	statusReady
)

type Model struct {
	config   *common.CloudConfig
	orgId    string
	asset    *cloudapiclient.AssetOutput
	ops      []cloudapiclient.OperationOutput
	fetchOps fetch.Model
	styles   uiCommon.Styles
	width    int
	height   int
	status   status
}

func GetAssetName(asset *cloudapiclient.AssetOutput) string {
	assetName := asset.CurrentAssetParameters.Data["name"]
	if assetName == "" {
		assetName = "N/A"
	}

	return assetName.(string)
}

func NewModel(config *common.CloudConfig, orgId string, asset *cloudapiclient.AssetOutput) *Model {
	m := &Model{
		orgId:  orgId,
		asset:  asset,
		styles: uiCommon.DefaultStyles(),
		status: statusInit,
		config: config,
	}

	return m
}

func Run(config *common.CloudConfig, orgId string, asset *cloudapiclient.AssetOutput) {
	p := tea.NewProgram(NewModel(config, orgId, asset), tea.WithAltScreen())
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
	return uiCommon.HelpView(items...)
}

func (m Model) bioView() string {
	s := m.styles.Logo.Render(GetAssetName(m.asset))
	s += "\n\n"
	s += uiCommon.KeyValueView(
		"Id", m.asset.Id,
		"Asset", m.asset.Asset,
	)
	s += m.opsTableView()
	return s
}

func (m Model) opsTableView() string {
	tbl := render.OpTable(m.ops)
	s := "\n\n\n"
	s += m.styles.Logo.Render("Operations")
	s += "  " + m.fetchOps.View()
	s += "\n"
	s += tbl.View()
	return s
}
