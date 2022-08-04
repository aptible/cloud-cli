package asset

import (
	"fmt"

	cloudapiclient "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/aptible/cloud-cli/config"
	libasset "github.com/aptible/cloud-cli/lib/asset"
	libenv "github.com/aptible/cloud-cli/lib/env"
	"github.com/aptible/cloud-cli/ui/fetch"
	"github.com/aptible/cloud-cli/ui/form"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// dsCreateRun - create a datastore
func vpcCreateRun() config.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := config.NewCloudConfig(viper.GetViper())
		org := config.Vconfig.GetString("org")
		env := config.Vconfig.GetString("env")

		formResult := form.FormResult{Org: org, Env: env}
		err := libenv.EnvForm(config, &formResult)
		if err != nil {
			return nil
		}

		name := args[0]

		vars := map[string]interface{}{
			"name": name,
		}
		params := cloudapiclient.AssetInput{
			Asset:           "aws__vpc__latest",
			AssetVersion:    "latest",
			AssetParameters: vars,
		}

		msg := fmt.Sprintf("creating vpc (%s)", name)
		model := fetch.NewModel(msg, func() (interface{}, error) {
			return config.Cc.CreateAsset(formResult.Org, formResult.Env, params)
		})

		result, err := fetch.WithOutput(model)
		if err != nil {
			return err
		}
		vpcTable := libasset.AssetTable(result.Result.(*cloudapiclient.AssetOutput))
		// TODO - print with tea
		fmt.Println("VPC(s) Created:")
		fmt.Println(vpcTable.View())

		return nil
	}
}

// dsDescribeRun - destroy datastore
func vpcDescribeRun() config.CobraRunE {
	return describeAsset()
}

// dsDestroyRun - destroy datastore
func vpcDestroyRun() config.CobraRunE {
	return destroyAsset()
}

// vpcListRun - list vpcs
func vpcListRun() config.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := config.NewCloudConfig(viper.GetViper())
		org := config.Vconfig.GetString("org")
		env := config.Vconfig.GetString("env")

		formResult := form.FormResult{Org: org, Env: env}
		err := libenv.EnvForm(config, &formResult)
		if err != nil {
			return nil
		}

		msg := fmt.Sprintf("getting vpcs with %+v", formResult)
		model := fetch.NewModel(msg, func() (interface{}, error) {
			return config.Cc.ListAssets(formResult.Org, formResult.Env)
		})

		rawResult, err := fetch.WithOutput(model)
		if err != nil {
			return err
		}
		if rawResult == nil {
			// TODO - print with tea
			fmt.Println("No vpcs found.")
			return nil
		}
		unfilteredResults := rawResult.Result.([]cloudapiclient.AssetOutput)
		filteredResults := libasset.FilterByType(unfilteredResults, []string{"vpc"})
		if len(filteredResults) == 0 {
			fmt.Println("No vpcs found.")
			return nil
		}

		vpcTable := libasset.AssetTable(filteredResults)
		// TODO - print with tea
		fmt.Println("VPC(s) List")
		fmt.Println(vpcTable.View())

		return nil
	}
}

func NewVPCCmd() *cobra.Command {
	vpcCmd := &cobra.Command{
		Use:     "network",
		Short:   "the network subcommand helps manage your Aptible network assets.",
		Long:    `The network subcommand helps manage your Aptible network assets.`,
		Aliases: []string{"v", "vpc"},
	}

	vpcCreateCmd := &cobra.Command{
		Use:     "create [asset_name]",
		Short:   "provision a new network.",
		Long:    `The network create command will provision a new network.`,
		Aliases: []string{"c", "deploy"},
		Args:    cobra.ExactArgs(1),
		RunE:    vpcCreateRun(),
	}

	vpcDestroyCmd := &cobra.Command{
		Use:     "destroy [asset_id]",
		Short:   "permanently remove the network.",
		Long:    `The network destroy command will permanently remove the network.`,
		Aliases: []string{"d", "delete", "rm", "remove"},
		Args:    cobra.ExactArgs(1),
		RunE:    vpcDestroyRun(),
	}

	vpcListCmd := &cobra.Command{
		Use:     "list",
		Short:   "list all networks within an organization.",
		Long:    `The network list command will list all networks within an organization.`,
		Aliases: []string{"ls"},
		RunE:    vpcListRun(),
	}

	vpcDescribeCmd := &cobra.Command{
		Use:     "describe [asset_id]",
		Short:   "describe vpc",
		Long:    `The network describe command will provide more detail for a network`,
		Aliases: []string{"show"},
		RunE:    vpcDescribeRun(),
	}

	vpcCmd.AddCommand(vpcCreateCmd)
	vpcCmd.AddCommand(vpcDescribeCmd)
	vpcCmd.AddCommand(vpcDestroyCmd)
	vpcCmd.AddCommand(vpcListCmd)

	return vpcCmd
}
