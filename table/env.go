package table

import (
	cloudapiclient "github.com/aptible/cloud-api-clients/clients/go"
	uiCommon "github.com/aptible/cloud-cli/internal/ui/common"
	"github.com/evertras/bubble-table/table"
)

// prints out a table of environments
func EnvTable(orgOutput interface{}) table.Model {
	rows := make([]table.Row, 0)

	switch data := orgOutput.(type) {
	case []cloudapiclient.EnvironmentOutput:
		for _, env := range data {
			rows = append(rows, table.NewRow(table.RowData{
				"id":             env.Id,
				"name":           env.Name,
				"aws_account_id": *env.AwsAccountId,
			}))
		}
	case *cloudapiclient.EnvironmentOutput:
		rows = append(rows, table.NewRow(table.RowData{
			"id":             data.Id,
			"name":           data.Name,
			"aws_account_id": *data.AwsAccountId,
		}))
	}

	return table.New([]table.Column{
		table.NewColumn("id", "Environment Id", 40).WithStyle(uiCommon.DefaultRowStyle()),
		table.NewColumn("name", "Environment Name", 40).WithStyle(uiCommon.DefaultRowStyle()),
		table.NewColumn("aws_account_id", "AWS Account Id", 40).WithStyle(uiCommon.DefaultRowStyle()),
	}).WithRows(rows)
}
