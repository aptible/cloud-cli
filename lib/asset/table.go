package libasset

import (
	"strings"

	cac "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/aptible/cloud-cli/ui/common"
	"github.com/evertras/bubble-table/table"
)

// colorizeAssetFromStatus - common utility for assets to colorize rows in CLI based on asset status
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

// generateAssetRowFromData - generate a common table row for assets
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

// dataStoreTable - prints out a table of assets
func AssetTable(orgOutput interface{}) table.Model {
	rows := make([]table.Row, 0)

	switch data := orgOutput.(type) {
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
