package table

import (
	cloudapiclient "github.com/aptible/cloud-api-clients/clients/go"
	uiCommon "github.com/aptible/cloud-cli/internal/ui/common"
	"github.com/evertras/bubble-table/table"
)

// prints out a table of organizations
func OrgTable(orgOutput interface{}) table.Model {
	rows := make([]table.Row, 0)

	switch data := orgOutput.(type) {
	case []cloudapiclient.OrganizationOutput:
		for _, org := range data {
			rows = append(rows, table.NewRow(table.RowData{
				"id":     org.Id,
				"name":   org.Name,
				"aws_ou": *org.AwsOu,
			}))
		}
	case *cloudapiclient.OrganizationOutput:
		rows = append(rows, table.NewRow(table.RowData{
			"id":     data.Id,
			"name":   data.Name,
			"aws_ou": *data.AwsOu,
		}))
	}

	return table.New([]table.Column{
		table.NewColumn("id", "Organization Id", 40).WithStyle(uiCommon.DefaultRowStyle()),
		table.NewColumn("name", "Organization Name", 40).WithStyle(uiCommon.DefaultRowStyle()),
		table.NewColumn("aws_ou", "AWS OU", 40).WithStyle(uiCommon.DefaultRowStyle()),
	}).WithRows(rows)
}
