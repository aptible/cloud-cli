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
	"github.com/aptible/cloud-cli/internal/ui/form"
)

func generateVpcRowFromData(asset cloudapiclient.AssetOutput) table.Row {
	row := table.NewRow(table.RowData{
		"id":     asset.Id,
		"status": asset.Status,
		"name":   asset.CurrentAssetParameters.Data["name"],
	})
	return colorizeAssetFromStatus(asset, row)
}

// dataStoreTable - prints out a table of datastores
func vpcTable(orgOutput interface{}) table.Model {
	rows := make([]table.Row, 0)

	switch data := orgOutput.(type) {
	case []cloudapiclient.AssetOutput:
		for _, asset := range data {
			rows = append(rows, generateVpcRowFromData(asset))
		}
	case *cloudapiclient.AssetOutput:
		rows = append(rows, generateVpcRowFromData(*data))
	}

	return table.New([]table.Column{
		table.NewColumn("id", "Asset Id", 40).WithStyle(uiCommon.DefaultRowStyle()),
		table.NewColumn("status", "Status", 40).WithStyle(uiCommon.DefaultRowStyle()),
		table.NewColumn("name", "Name", 40).WithStyle(uiCommon.DefaultRowStyle()),
	}).WithRows(rows)
}

// dsCreateRun - create a datastore
func vpcCreateRun() common.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := common.NewCloudConfig(viper.GetViper())
		org := config.Vconfig.GetString("org")
		env := config.Vconfig.GetString("env")

		formResult := form.FormResult{Org: org, Env: env}
		err := form.EnvForm(config, &formResult)
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
		vpcTable := vpcTable(result.Result.(*cloudapiclient.AssetOutput))
		// TODO - print with tea
		fmt.Println("VPC(s) Created:")
		fmt.Println(vpcTable.View())

		return nil
	}
}

// dsDescribeRun - destroy datastore
func vpcDescribeRun() common.CobraRunE {
	return describeAsset()
}

// dsDestroyRun - destroy datastore
func vpcDestroyRun() common.CobraRunE {
	return destroyAsset()
}

// vpcListRun - list vpcs
func vpcListRun() common.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := common.NewCloudConfig(viper.GetViper())
		org := config.Vconfig.GetString("org")
		env := config.Vconfig.GetString("env")

		formResult := form.FormResult{Org: org, Env: env}
		err := form.EnvForm(config, &formResult)
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
		filteredResults := common.FilterAssetsByType(unfilteredResults, []string{"vpc"})
		if len(filteredResults) == 0 {
			fmt.Println("No vpcs found.")
			return nil
		}

		vpcTable := vpcTable(filteredResults)
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
