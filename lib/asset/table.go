package libasset

import (
	"fmt"
	"os"
	"strings"

	cac "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/aptible/cloud-cli/ui/common"
	"github.com/evertras/bubble-table/table"
	"golang.org/x/term"
)

func colorizeFromStatus(asset cac.AssetOutput, row table.Row) table.Row {
	switch asset.Status {
	case cac.ASSETSTATUS_DEPLOYED:
		return row.WithStyle(common.ActiveRowStyle())
	case cac.ASSETSTATUS_DEPLOYING,
		cac.ASSETSTATUS_PENDING,
		cac.ASSETSTATUS_DESTROYING,
		cac.ASSETSTATUS_REQUESTED:
		return row.WithStyle(common.PendingRowStyle())
	case cac.ASSETSTATUS_DESTROYED:
		return row.WithStyle(common.DisabledRowStyle())
	default:
		return row.WithStyle(common.DefaultRowStyle())
	}
}

func generateRowFromData(asset cac.AssetOutput) table.Row {
	assetName := GetName(asset)
	assetStr := strings.Split(asset.Asset, "__")
	row := table.NewRow(table.RowData{
		"id":            asset.Id,
		"status":        asset.Status,
		"name":          assetName,
		"cloud":         assetStr[0],
		"asset_type":    assetStr[1],
		"asset_version": assetStr[2],
	})
	return colorizeFromStatus(asset, row)
}

func AssetTable(output interface{}) table.Model {
	rows := make([]table.Row, 0)

	switch data := output.(type) {
	case []cac.AssetOutput:
		for _, asset := range data {
			rows = append(rows, generateRowFromData(asset))
		}
	case *cac.AssetOutput:
		rows = append(rows, generateRowFromData(*data))
	}

	return table.New([]table.Column{
		table.NewColumn("id", "Id", 40).WithStyle(common.DefaultRowStyle()),
		table.NewColumn("status", "Status", 40).WithStyle(common.DefaultRowStyle()),
		table.NewColumn("name", "Name", 20).WithStyle(common.DefaultRowStyle()),
		table.NewColumn("cloud", "Cloud", 20).WithStyle(common.DefaultRowStyle()),
		table.NewColumn("asset_type", "Type", 20).WithStyle(common.DefaultRowStyle()),
		table.NewColumn("asset_version", "Version", 20).WithStyle(common.DefaultRowStyle()),
	}).WithRows(rows)
}

func generateRowFromBundleData(bundle cac.AssetBundle) table.Row {
	row := table.NewRow(table.RowData{
		"id":          bundle.Identifier,
		"name":        bundle.Name,
		"description": bundle.Description,
	})
	return row
}

func AssetBundleTable(output interface{}) table.Model {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println(err)
	}
	rows := make([]table.Row, 0)

	switch data := output.(type) {
	case []cac.AssetBundle:
		for _, bundle := range data {
			rows = append(rows, generateRowFromBundleData(bundle))
		}
	case *cac.AssetBundle:
		rows = append(rows, generateRowFromBundleData(*data))
	}

	return table.New([]table.Column{
		table.NewFlexColumn("id", "Id", 1).WithStyle(common.LeftRowStyle()),
		table.NewFlexColumn("name", "Name", 1).WithStyle(common.LeftRowStyle()),
		table.NewFlexColumn("description", "Description", 3).WithStyle(common.LeftRowStyle()),
	}).WithRows(rows).WithTargetWidth(width)
}
