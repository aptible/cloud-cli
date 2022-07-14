package assets

import (
	"fmt"

	cloudapiclient "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/evertras/bubble-table/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aptible/cloud-cli/internal/common"
	uiCommon "github.com/aptible/cloud-cli/internal/ui/common"
	"github.com/aptible/cloud-cli/internal/ui/fetch"
)

func colorizeAssetFromStatus(asset cloudapiclient.AssetOutput, row table.Row) table.Row {
	switch asset.Status {
	case "DEPLOYED":
		return row.WithStyle(uiCommon.ActiveRowStyle())
	case "PENDING", "DEPLOYING", "DESTROYING":
		return row.WithStyle(uiCommon.PendingRowStyle())
	case "DESTROYED":
		return row.WithStyle(uiCommon.DisabledRowStyle())
	default:
		return row.WithStyle(uiCommon.DefaultRowStyle())
	}
	return row
}

func destroyAsset(_ *cobra.Command, args []string) error {
	config := common.NewCloudConfig(viper.GetViper())
	orgId := config.Vconfig.GetString("org")
	assetId := args[0]

	if env == "" {
		return fmt.Errorf("must provide env")
	}

	msg := fmt.Sprintf("destroying asset %s (v%s)", engine, engineVersion)
	model := fetch.NewModel(msg, func() (interface{}, int, error) {
		status, err := config.Cc.DestroyAsset(orgId, env, assetId)
		return nil, status, err
	})
	_, err := fetch.WithOutput(model)
	if err != nil {
		return err
	}

	fmt.Printf("destroying asset with id: %+v\n", assetId)
	return nil
}
