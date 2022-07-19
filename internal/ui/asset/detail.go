package asset

import (
	"fmt"
	"log"
	"time"

	cloudapiclient "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/aptible/cloud-cli/internal/ui/common"
	render "github.com/aptible/cloud-cli/table"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	asset  *cloudapiclient.AssetOutput
	ops    []cloudapiclient.OperationOutput
	styles common.Styles
	width  int
	height int
}

type tickMsg time.Time

func GetAssetName(asset *cloudapiclient.AssetOutput) string {
	assetName := asset.CurrentAssetParameters.Data["name"]
	if assetName == "" {
		assetName = "N/A"
	}

	return assetName.(string)
}

func NewModel(asset *cloudapiclient.AssetOutput, ops []cloudapiclient.OperationOutput) *Model {
	m := &Model{
		asset:  asset,
		ops:    ops,
		styles: common.DefaultStyles(),
	}

	return m
}

func Run(asset *cloudapiclient.AssetOutput, ops []cloudapiclient.OperationOutput) {
	p := tea.NewProgram(NewModel(asset, ops), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.styles.Window.Height(m.height - 4)
		m.styles.Window.Width(m.width - 2)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		default:
			return m, nil
		}
	}

	return m, nil
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
	s := m.styles.Logo.Render(GetAssetName(m.asset))
	s += "\n\n"
	s += common.KeyValueView(
		"Id", m.asset.Id,
		"Asset", m.asset.Asset,
	)
	s += m.opsTableView()
	return s
}

func (m Model) opsTableView() string {
	tbl := render.OperationTable(m.ops)
	s := "\n\n\n"
	s += m.styles.Logo.Render("Operations")
	s += "\n"
	s += tbl.View()
	return s
}
