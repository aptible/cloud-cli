package asset

import (
	"fmt"
	"strings"

	cac "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aptible/cloud-cli/config"
	"github.com/aptible/cloud-cli/lib/asset"
	libenv "github.com/aptible/cloud-cli/lib/env"
	"github.com/aptible/cloud-cli/ui/fetch"
	"github.com/aptible/cloud-cli/ui/form"
)

type AssetOptions struct {
	AssetName     string
	AssetType     string
	VpcName       string
	Engine        string
	EngineVersion string
}

var assetOptions = AssetOptions{}

// describeAsset - aliased func but describes any given asset by its asset id, env id (rds/vpc for example use this)
func describeAsset() config.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := config.NewCloudConfig(viper.GetViper())
		org := config.Vconfig.GetString("org")
		env := config.Vconfig.GetString("env")
		assetId := args[0]

		formResult := form.FormResult{Org: org, Env: env}
		err := libenv.EnvForm(config, &formResult)
		if err != nil {
			return nil
		}

		msg := fmt.Sprintf("describing asset %s", assetId)
		model := fetch.NewModel(msg, func() (interface{}, error) {
			return config.Cc.DescribeAsset(formResult.Org, formResult.Env, assetId)
		})
		data, err := fetch.WithOutput(model)
		if err != nil {
			return err
		}
		asset := data.Result.(*cac.AssetOutput)

		libasset.RunDetail(config, formResult.Org, asset)

		return nil
	}
}

func assetDescribeRun() config.CobraRunE {
	return describeAsset()
}

// destroyAsset - aliased func but also can destroy assets on top level (rds/vpc for example use this)
func destroyAsset() config.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := config.NewCloudConfig(viper.GetViper())
		org := config.Vconfig.GetString("org")
		env := config.Vconfig.GetString("env")
		assetId := args[0]

		formResult := form.FormResult{Org: org, Env: env}
		err := libenv.EnvForm(config, &formResult)
		if err != nil {
			return nil
		}

		msg := fmt.Sprintf("destroying asset %s (v%s)", engine, engineVersion)
		model := fetch.NewModel(msg, func() (interface{}, error) {
			err := config.Cc.DestroyAsset(formResult.Org, formResult.Env, assetId)
			return nil, err
		})
		_, err = fetch.WithOutput(model)
		if err != nil {
			return err
		}

		fmt.Printf("Started request to destroy asset with id: %+v\n", assetId)
		return nil
	}
}

// assetsCreateRun - create an asset
func assetsCreateRun() config.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := config.NewCloudConfig(viper.GetViper())
		org := config.Vconfig.GetString("org")
		env := config.Vconfig.GetString("env")

		formResult := form.FormResult{
			Org:           org,
			Env:           env,
			AssetType:     assetOptions.AssetType,
			AssetName:     assetOptions.AssetName,
			Engine:        assetOptions.Engine,
			EngineVersion: assetOptions.EngineVersion,
		}
		err := libasset.AssetCreateForm(config, &formResult)
		if err != nil {
			return err
		}

		vars := map[string]interface{}{
			"name":           formResult.AssetName,
			"engine":         formResult.Engine,
			"engine_version": formResult.EngineVersion,
			"vpc_name":       formResult.VpcName,
		}
		// TODO: asset type => what's required here is not the same
		_type := strings.Replace(formResult.AssetType, "aws/", "", 1)
		params := cac.AssetInput{
			Asset:           fmt.Sprintf("aws__%s__latest", _type),
			AssetVersion:    "latest",
			AssetParameters: vars,
		}

		msg := fmt.Sprintf("creating asset %s (v%s)", formResult.Engine, formResult.EngineVersion)
		model := fetch.NewModel(msg, func() (interface{}, error) {
			return config.Cc.CreateAsset(formResult.Org, formResult.Env, params)
		})

		result, err := fetch.WithOutput(model)
		if err != nil {
			return err
		}
		res := result.Result.(*cac.AssetOutput)

		fmt.Printf("Result: %+v\n", res)
		return nil
	}
}

// assetsDestroyRun - destory an asset
func assetsDestroyRun() config.CobraRunE {
	return destroyAsset()
}

// assetsListRun - list all possible assets with config fields
func assetsListRun() config.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := config.NewCloudConfig(viper.GetViper())
		org := config.Vconfig.GetString("org")
		env := config.Vconfig.GetString("env")

		formResult := form.FormResult{Org: org, Env: env}
		err := libenv.EnvForm(config, &formResult)
		if err != nil {
			return err
		}

		msg := fmt.Sprintf("getting assets with %+v", formResult)
		model := fetch.NewModel(msg, func() (interface{}, error) {
			return config.Cc.ListAssets(formResult.Org, formResult.Env)
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

		results := rawResult.Result.([]cac.AssetOutput)
		if len(results) == 0 {
			// TODO - print with tea
			fmt.Println("No assets found.")
			return nil
		}

		dsTable := libasset.AssetTable(results)
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
		Use:     "create [asset_type] [name]",
		Short:   "provision a new asset.",
		Long:    `The asset create command will provision a new asset.`,
		Aliases: []string{"c", "deploy"},
		RunE:    assetsCreateRun(),
	}

	assetDestroyCmd := &cobra.Command{
		Use:     "destroy [asset_id]",
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

	assetDescribeCmd := &cobra.Command{
		Use:     "describe",
		Short:   "Show asset detail",
		Long:    `The assets describe command will provide more detail about the asset`,
		Aliases: []string{"show"},
		RunE:    assetDescribeRun(),
	}

	assetCreateCmd.Flags().StringVarP(&assetOptions.VpcName, "vpc-name", "", "", "vpc name to create the asset in")
	assetCreateCmd.Flags().StringVarP(&assetOptions.AssetName, "asset-name", "", "", "asset name")
	assetCreateCmd.Flags().StringVarP(&assetOptions.AssetType, "asset-type", "", "", "asset type")
	assetCreateCmd.Flags().StringVarP(&assetOptions.Engine, "engine", "", "", "engine")
	assetCreateCmd.Flags().StringVarP(&assetOptions.EngineVersion, "engine-version", "", "", "engine version")

	assetCmd.AddCommand(assetCreateCmd)
	assetCmd.AddCommand(assetDestroyCmd)
	assetCmd.AddCommand(assetListCmd)
	assetCmd.AddCommand(assetDescribeCmd)

	return assetCmd
}
