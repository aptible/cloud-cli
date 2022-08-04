package libop

import (
	cac "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/aptible/cloud-cli/ui/common"
	"github.com/evertras/bubble-table/table"
)

// colorizeOperationFromStatus - common utility for assets to colorize rows in CLI based on asset status
func colorizeOperationFromStatus(operation cac.OperationOutput, row table.Row) table.Row {
	switch *operation.Status.Get() {
	case cac.OPERATIONSTATUS_COMPLETE:
		return row.WithStyle(common.ActiveRowStyle())
	case cac.OPERATIONSTATUS_IN_PROGRESS,
		cac.OPERATIONSTATUS_PAUSED,
		cac.OPERATIONSTATUS_PENDING:
		return row.WithStyle(common.PendingRowStyle())
	case cac.OPERATIONSTATUS_CANCELED, cac.OPERATIONSTATUS_FAILED:
		return row.WithStyle(common.DisabledRowStyle())
	default:
		return row.WithStyle(common.DefaultRowStyle())
	}
}

// generateAssetRowFromData - generate a common table row for assets
func generateOpRowFromData(op cac.OperationOutput) table.Row {
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
	case []cac.OperationOutput:
		for _, op := range data {
			rows = append(rows, generateOpRowFromData(op))
		}
	case *cac.OperationOutput:
		rows = append(rows, generateOpRowFromData(*data))
	}

	return table.New([]table.Column{
		table.NewColumn("id", "Id", 40).WithStyle(common.DefaultRowStyle()),
		table.NewColumn("type", "Type", 20).WithStyle(common.DefaultRowStyle()),
		table.NewColumn("status", "Status", 40).WithStyle(common.DefaultRowStyle()),
	}).WithRows(rows)
}
