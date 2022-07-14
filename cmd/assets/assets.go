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

// colorizeAssetFromStatus - common utility for assets to colorize rows in CLI based on asset status
func colorizeAssetFromStatus(asset cloudapiclient.AssetOutput, row table.Row) table.Row {
	switch asset.Status {
	case "DEPLOYED":
		return row.WithStyle(uiCommon.ActiveRowStyle())
	case "PENDING", "DEPLOYING", "DESTROYING", "REQUESTED":
		return row.WithStyle(uiCommon.PendingRowStyle())
	case "DESTROYED":
		return row.WithStyle(uiCommon.DisabledRowStyle())
	default:
		return row.WithStyle(uiCommon.DefaultRowStyle())
	}
}

// listAssets - list all possible assets with common fields
func listAssets() common.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		return nil
	}
}

// createAsset - entry point to create an asset barebones
func createAsset() common.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := common.NewCloudConfig(viper.GetViper())
		orgID := config.Vconfig.GetString("org")
		envID := args[0]
		assetType := args[1]
		name := args[2]

		if engine == "" {
			return fmt.Errorf("must provide engine")
		}
		if engineVersion == "" {
			return fmt.Errorf("must provide engine version")
		}

		vars := map[string]interface{}{
			"name": name,
		}
		params := cloudapiclient.AssetInput{
			Asset:           fmt.Sprintf("aws__%s__latest", assetType),
			AssetVersion:    "latest",
			AssetParameters: vars,
		}

		msg := fmt.Sprintf("creating asset %s (v%s)", engine, engineVersion)
		model := fetch.NewModel(msg, func() (interface{}, int, error) {
			return config.Cc.CreateAsset(orgID, envID, params)
		})

		result, err := fetch.WithOutput(model)
		if err != nil {
			return err
		}
		res := result.Result.(*cloudapiclient.AssetOutput)

		fmt.Printf("Result: %+v\n", res)
		return nil
	}
}

// describeAsset - commonly aliased func but describes any given asset by its asset id, env id (rds/vpc for example use this)
func describeAsset() common.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		return nil
	}
}

// destroyAsset - commonly aliased func but also can destroy assets on top level (rds/vpc for example use this)
func destroyAsset() common.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
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

		fmt.Printf("Started request to destroy asset with id: %+v\n", assetId)
		return nil
	}
}
