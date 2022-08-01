package table

import (
	cloudapiclient "github.com/aptible/cloud-api-clients/clients/go"
	uiCommon "github.com/aptible/cloud-cli/internal/ui/common"
	"github.com/evertras/bubble-table/table"
)

// colorizeOperationFromStatus - common utility for assets to colorize rows in CLI based on asset status
func colorizeOperationFromStatus(operation cloudapiclient.OperationOutput, row table.Row) table.Row {
	switch *operation.Status.Get() {
	case cloudapiclient.OPERATIONSTATUS_COMPLETE:
		return row.WithStyle(uiCommon.ActiveRowStyle())
	case cloudapiclient.OPERATIONSTATUS_IN_PROGRESS,
		cloudapiclient.OPERATIONSTATUS_PAUSED,
		cloudapiclient.OPERATIONSTATUS_PENDING:
		return row.WithStyle(uiCommon.PendingRowStyle())
	case cloudapiclient.OPERATIONSTATUS_CANCELED, cloudapiclient.OPERATIONSTATUS_FAILED:
		return row.WithStyle(uiCommon.DisabledRowStyle())
	default:
		return row.WithStyle(uiCommon.DefaultRowStyle())
	}
}

// generateAssetRowFromData - generate a common table row for assets
func generateOpRowFromData(op cloudapiclient.OperationOutput) table.Row {
	row := table.NewRow(table.RowData{
		"id":     op.Id,
		"type":   *op.OperationType.Get(),
		"status": *op.Status.Get(),
	})
	return colorizeOperationFromStatus(op, row)
}

// dataStoreTable - prints out a table of operations
func OpTable(orgOutput interface{}) table.Model {
	rows := make([]table.Row, 0)

	switch data := orgOutput.(type) {
	case []cloudapiclient.OperationOutput:
		for _, op := range data {
			rows = append(rows, generateOpRowFromData(op))
		}
	case *cloudapiclient.OperationOutput:
		rows = append(rows, generateOpRowFromData(*data))
	}

	return table.New([]table.Column{
		table.NewColumn("id", "Id", 40).WithStyle(uiCommon.DefaultRowStyle()),
		table.NewColumn("type", "Type", 20).WithStyle(uiCommon.DefaultRowStyle()),
		table.NewColumn("status", "Status", 40).WithStyle(uiCommon.DefaultRowStyle()),
	}).WithRows(rows)
}
