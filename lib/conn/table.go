package libconn

import (
	"fmt"

	cac "github.com/aptible/cloud-api-clients/clients/go"
	libasset "github.com/aptible/cloud-cli/lib/asset"
	"github.com/aptible/cloud-cli/ui/common"
	"github.com/evertras/bubble-table/table"
)

// prints out a table of connections
func ConnTable(connOutput interface{}) table.Model {
	rows := make([]table.Row, 0)

	switch data := connOutput.(type) {
	case []cac.ConnectionOutput:
		for _, conn := range data {
			inc := ""
			if conn.HasIncomingConnectionAsset() {
				inc = libasset.GetName(*conn.IncomingConnectionAsset)
			}
			out := ""
			if conn.HasOutgoingConnectionAsset() {
				out = libasset.GetName(*conn.OutgoingConnectionAsset)
			}
			rows = append(rows, table.NewRow(table.RowData{
				"id":   conn.Id,
				"conn": fmt.Sprintf("%s => %s", out, inc),
			}))
		}
	case *cac.ConnectionOutput:
		inc := ""
		out := ""
		if data.HasIncomingConnectionAsset() {
			inc = libasset.GetName(*data.IncomingConnectionAsset)
		}
		if data.HasOutgoingConnectionAsset() {
			out = libasset.GetName(*data.OutgoingConnectionAsset)
		}
		rows = append(rows, table.NewRow(table.RowData{
			"id":   data.Id,
			"conn": fmt.Sprintf("%s => %s", out, inc),
		}))
	}

	return table.New([]table.Column{
		table.NewColumn("id", "Id", 40).WithStyle(common.DefaultRowStyle()),
		table.NewColumn("conn", "Connection", 40).WithStyle(common.DefaultRowStyle()),
	}).WithRows(rows)
}
