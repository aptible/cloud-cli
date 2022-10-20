package libenv

import (
	cac "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/aptible/cloud-cli/ui/common"
	"github.com/evertras/bubble-table/table"
)

func safeString(str *string) string {
	if str == nil {
		return "unknown"
	}

	return *str
}

// prints out a table of environments
func EnvTable(orgOutput interface{}) table.Model {
	rows := make([]table.Row, 0)

	switch data := orgOutput.(type) {
	case []cac.EnvironmentOutput:
		for _, env := range data {
			rows = append(rows, table.NewRow(table.RowData{
				"id":             env.Id,
				"name":           env.Name,
				"aws_account_id": safeString(env.AwsAccountId),
			}))
		}
	case *cac.EnvironmentOutput:
		rows = append(rows, table.NewRow(table.RowData{
			"id":             data.Id,
			"name":           data.Name,
			"aws_account_id": safeString(data.AwsAccountId),
		}))
	}

	return table.New([]table.Column{
		table.NewColumn("id", "Environment Id", 40).WithStyle(common.DefaultRowStyle()),
		table.NewColumn("name", "Environment Name", 40).WithStyle(common.DefaultRowStyle()),
		table.NewColumn("aws_account_id", "AWS Account Id", 40).WithStyle(common.DefaultRowStyle()),
	}).WithRows(rows)
}
