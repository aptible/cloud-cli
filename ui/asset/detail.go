package assetui

import (
	"fmt"
	"log"
	"reflect"
	"time"

	cac "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/aptible/cloud-cli/config"
	"github.com/aptible/cloud-cli/lib/asset"
	"github.com/aptible/cloud-cli/lib/conn"
	"github.com/aptible/cloud-cli/lib/op"
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
	asset    *cac.AssetOutput
	ops      []cac.OperationOutput
	fetchOps fetch.Model
	styles   common.Styles
	width    int
	height   int
	status   status
}

func NewDetailModel(config *config.CloudConfig, orgId string, asset *cac.AssetOutput) *Model {
	m := &Model{
		orgId:  orgId,
		asset:  asset,
		styles: common.DefaultStyles(),
		status: statusInit,
		config: config,
	}

	return m
}

func RunDetail(config *config.CloudConfig, orgId string, asset *cac.AssetOutput) {
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
		m.ops = msg.Result.([]cac.OperationOutput)
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

func infToStr(inf interface{}) string {
	if inf != nil {
		infType := reflect.TypeOf(inf).String()
		if infType == "string" || infType == "int" {
			return inf.(string)
		}
	}
	return ""
}

func (m Model) bioView() string {
	vs := []string{
		"Id", m.asset.Id,
		"Asset", m.asset.Asset,
		"VPC", infToStr(m.asset.CurrentAssetParameters.Data["vpc_name"]),
		"Engine", infToStr(m.asset.CurrentAssetParameters.Data["engine"]),
		"Engine Version", infToStr(m.asset.CurrentAssetParameters.Data["engine_version"]),
	}
	s := m.styles.Logo.Render(libasset.GetName(*m.asset))
	s += "\n\n"
	s += common.KeyValueView(vs...)
	s += m.opsTableView()
	s += m.connTableView()
	return s
}

func (m Model) connTableView() string {
	if len(m.asset.Connections) == 0 {
		return ""
	}

	tbl := libconn.ConnTable(m.asset.Connections)
	s := "\n\n\n"
	s += m.styles.Logo.Render("Connections")
	s += "\n"
	s += tbl.View()
	return s
}

func (m Model) opsTableView() string {
	if len(m.ops) == 0 {
		return ""
	}

	tbl := libop.OpTable(m.ops)
	s := "\n\n\n"
	s += m.styles.Logo.Render("Operations")
	s += "  " + m.fetchOps.View()
	s += "\n"
	s += tbl.View()
	return s
}
