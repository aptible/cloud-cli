package cmd

import (
	"fmt"

	apiclient "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/aptible/cloud-cli/ui/fetch"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func assetCreateRun() RunE {
	return func(cmd *cobra.Command, args []string) error {
		config := NewCloudConfig(viper.GetViper())
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
		params := apiclient.AssetInput{
			Asset:           fmt.Sprintf("aws__%s__latest", assetType),
			AssetVersion:    "latest",
			AssetParameters: vars,
		}

		msg := fmt.Sprintf("creating asset %s (v%s)", engine, engineVersion)
		model := fetch.NewModel(msg, func() (interface{}, error) {
			return config.Cc.CreateAsset(orgID, envID, params)
		})

		result, err := fetch.FetchWithOutput(model)
		if err != nil {
			return err
		}
		res := result.Result.(*apiclient.AssetOutput)

		fmt.Printf("Result: %+v\n", res)
		return nil
	}
}

func assetDestroyRun() RunE {
	return func(cmd *cobra.Command, args []string) error {
		config := NewCloudConfig(viper.GetViper())
		orgID := config.Vconfig.GetString("org")
		envID := args[0]
		assetID := args[1]

		model := fetch.NewModel("destroying asset", func() (interface{}, error) {
			err := config.Cc.DestroyAsset(orgID, envID, assetID)
			return nil, err
		})

		err := fetch.FetchAny(model)
		return err
	}
}

func assetListRun() RunE {
	return func(cmd *cobra.Command, args []string) error {
		config := NewCloudConfig(viper.GetViper())
		orgID := config.Vconfig.GetString("org")
		envID := args[0]

		model := fetch.NewModel("fetching assets", func() (interface{}, error) {
			return config.Cc.ListAssets(orgID, envID)
		})
		result, err := fetch.FetchWithOutput(model)
		if err != nil {
			return err
		}

		assets := result.Result.([]apiclient.AssetOutput)

		for _, asset := range assets {
			fmt.Printf("%s %s\n", asset.Id, asset.Asset)
		}
		return nil
	}
}

func NewAssetCmd() *cobra.Command {
	assetCmd := &cobra.Command{
		Use:     "asset",
		Short:   "The asset subcommand helps manage your Aptible assets.",
		Long:    `The asset subcommand helps manage your Aptible assets.`,
		Aliases: []string{"ass", "a"},
	}

	assetCreateCmd := &cobra.Command{
		Use:     "create",
		Short:   "create a new asset.",
		Long:    `The asset create command will create a new asset.`,
		Aliases: []string{"c"},
		Args:    cobra.MinimumNArgs(3),
		RunE:    assetCreateRun(),
	}

	assetDestroyCmd := &cobra.Command{
		Use:     "destroy",
		Short:   "permentantly remove the asset.",
		Long:    `The destroy command will permentantly remove the asset.`,
		Aliases: []string{"d", "delete", "rm", "remove"},
		Args:    cobra.MinimumNArgs(2),
		RunE:    assetDestroyRun(),
	}

	assetListCmd := &cobra.Command{
		Use:     "list",
		Short:   "list all assets within an organization.",
		Long:    `The asset list command will list all assets within an organization.`,
		Aliases: []string{"ls"},
		Args:    cobra.MinimumNArgs(1),
		RunE:    assetListRun(),
	}

	assetCmd.AddCommand(assetCreateCmd)
	assetCmd.AddCommand(assetDestroyCmd)
	assetCmd.AddCommand(assetListCmd)

	return assetCmd
}
