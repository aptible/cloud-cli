package table

import (
	"strings"

	cloudapiclient "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/evertras/bubble-table/table"

	"github.com/aptible/cloud-cli/internal/common"
	uiCommon "github.com/aptible/cloud-cli/internal/ui/common"
)

// colorizeAssetFromStatus - common utility for assets to colorize rows in CLI based on asset status
func colorizeAssetFromStatus(asset cloudapiclient.AssetOutput, row table.Row) table.Row {
	switch asset.Status {
	case cloudapiclient.ASSETSTATUS_DEPLOYED:
		return row.WithStyle(uiCommon.ActiveRowStyle())
	case cloudapiclient.ASSETSTATUS_DEPLOYING,
		cloudapiclient.ASSETSTATUS_PENDING,
		cloudapiclient.ASSETSTATUS_DESTROYING,
		cloudapiclient.ASSETSTATUS_REQUESTED:
		return row.WithStyle(uiCommon.PendingRowStyle())
	case cloudapiclient.ASSETSTATUS_DESTROYED:
		return row.WithStyle(uiCommon.DisabledRowStyle())
	default:
		return row.WithStyle(uiCommon.DefaultRowStyle())
	}
}

// generateAssetRowFromData - generate a common table row for assets
func generateAssetRowFromData(asset cloudapiclient.AssetOutput) table.Row {
	assetName := common.GetAssetName(asset)
	assetStr := strings.Split(asset.Asset, "__")
	row := table.NewRow(table.RowData{
		"id":            asset.Id,
		"status":        asset.Status,
		"name":          assetName,
		"cloud":         assetStr[0],
		"asset_type":    assetStr[1],
		"asset_version": assetStr[2],
	})
	return colorizeAssetFromStatus(asset, row)
}

// dataStoreTable - prints out a table of assets
func AssetTable(orgOutput interface{}) table.Model {
	rows := make([]table.Row, 0)

	switch data := orgOutput.(type) {
	case []cloudapiclient.AssetOutput:
		for _, asset := range data {
			rows = append(rows, generateAssetRowFromData(asset))
		}
	case *cloudapiclient.AssetOutput:
		rows = append(rows, generateAssetRowFromData(*data))
	}

	return table.New([]table.Column{
		table.NewColumn("id", "Id", 40).WithStyle(uiCommon.DefaultRowStyle()),
		table.NewColumn("status", "Status", 40).WithStyle(uiCommon.DefaultRowStyle()),
		table.NewColumn("name", "Name", 20).WithStyle(uiCommon.DefaultRowStyle()),
		table.NewColumn("cloud", "Cloud", 20).WithStyle(uiCommon.DefaultRowStyle()),
		table.NewColumn("asset_type", "Type", 20).WithStyle(uiCommon.DefaultRowStyle()),
		table.NewColumn("asset_version", "Version", 20).WithStyle(uiCommon.DefaultRowStyle()),
	}).WithRows(rows)
}
