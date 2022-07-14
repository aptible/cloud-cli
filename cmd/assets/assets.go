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

// generateAssetRowFromData - generate a common table row for assets
func generateAssetRowFromData(asset cloudapiclient.AssetOutput) table.Row {
	assetName := asset.CurrentAssetParameters.Data["name"]
	if assetName == "" {
		assetName = "N/A"
	}
	row := table.NewRow(table.RowData{
		"id":     asset.Id,
		"status": asset.Status,
		"name":   assetName,
	})
	return colorizeAssetFromStatus(asset, row)
}

// dataStoreTable - prints out a table of assets
func assetTable(orgOutput interface{}) table.Model {
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
		envId := config.Vconfig.GetString("env")
		assetId := args[0]

		if envId == "" {
			return fmt.Errorf("must provide env")
		}

		msg := fmt.Sprintf("destroying asset %s (v%s)", engine, engineVersion)
		model := fetch.NewModel(msg, func() (interface{}, int, error) {
			status, err := config.Cc.DestroyAsset(orgId, envId, assetId)
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

// assetsCreateRun - create an asset
func assetsCreateRun() common.CobraRunE {
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

// assetsDestroyRun - destory an asset
func assetsDestroyRun() common.CobraRunE {
	return destroyAsset()
}

// assetsListRun - list all possible assets with common fields
func assetsListRun() common.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := common.NewCloudConfig(viper.GetViper())
		orgId := config.Vconfig.GetString("org")
		envId := config.Vconfig.GetString("env")

		msg := fmt.Sprintf("getting datastores with env id: %s and org id: %s", envId, orgId)
		model := fetch.NewModel(msg, func() (interface{}, int, error) {
			return config.Cc.ListAssets(orgId, envId)
		})

		rawResult, err := fetch.WithOutput(model)
		if err != nil {
			return err
		}
		if rawResult == nil {
			// TODO - print with tea
			fmt.Println("No datastores found.")
			return nil
		}

		results := rawResult.Result.([]cloudapiclient.AssetOutput)
		if len(results) == 0 {
			// TODO - print with tea
			fmt.Println("No assets found.")
			return nil
		}

		dsTable := assetTable(results)
		// TODO - print with tea
		fmt.Println("Asset(s) List")
		fmt.Println(dsTable.View())

		return nil
	}
}

// NewAssetCmd - generate a new asset
func NewAssetCmd() *cobra.Command {
	assetCmd := &cobra.Command{
		Use:     "asset",
		Short:   "the asset subcommand helps manage your Aptible assets.",
		Long:    `The asset subcommand helps manage your Aptible assets.`,
		Aliases: []string{"a"},
	}

	assetCreateCmd := &cobra.Command{
		Use:     "create",
		Short:   "provision a new asset.",
		Long:    `The asset create command will provision a new asset.`,
		Aliases: []string{"c", "deploy"},
		RunE:    assetsCreateRun(),
	}

	assetDestroyCmd := &cobra.Command{
		Use:     "destroy",
		Short:   "permanently remove the asset.",
		Long:    `The asset destroy command will permanently remove the asset.`,
		Aliases: []string{"d", "delete", "rm", "remove"},
		Args:    cobra.MinimumNArgs(1),
		RunE:    assetsDestroyRun(),
	}

	assetListCmd := &cobra.Command{
		Use:     "list",
		Short:   "list all assets within an organization.",
		Long:    `The assets list command will list all assets within an organization.`,
		Aliases: []string{"ls"},
		RunE:    assetsListRun(),
	}

	assetCreateCmd.Flags().StringVarP(&vpcName, "vpc-name", "", "", "asset variables map")

	assetCmd.AddCommand(assetCreateCmd)
	assetCmd.AddCommand(assetDestroyCmd)
	assetCmd.AddCommand(assetListCmd)

	return assetCmd
}
