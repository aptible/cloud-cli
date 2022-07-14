package assets

import (
	"fmt"
	"strings"

	cloudapiclient "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/evertras/bubble-table/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aptible/cloud-cli/internal/common"
	uiCommon "github.com/aptible/cloud-cli/internal/ui/common"
	"github.com/aptible/cloud-cli/internal/ui/fetch"
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
		orgId := config.Vconfig.GetString("org")
		envId := args[0]
		name := args[1]

		vars := map[string]interface{}{
			"name": name,
		}
		params := cloudapiclient.AssetInput{
			Asset:           "aws__vpc__latest",
			AssetVersion:    "latest",
			AssetParameters: vars,
		}

		msg := fmt.Sprintf("creating vpc %s (v%s)", engine, engineVersion)
		model := fetch.NewModel(msg, func() (interface{}, int, error) {
			return config.Cc.CreateAsset(orgId, envId, params)
		})

		result, err := fetch.WithOutput(model)
		if err != nil {
			return err
		}
		vpcTable := vpcTable(result.Result.(*cloudapiclient.AssetOutput))
		// TODO - print with tea
		fmt.Println("VPC(s) List")
		fmt.Println(vpcTable.View())

		return nil
	}
}

// dsDestroyRun - destroy datastore
func vpcDestroyRun() common.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		fmt.Println(fmt.Sprintf("Destroying vpc id: %s", args[0]))
		return destroyAsset(cmd, args)
	}
}

// vpcListRun - list vpcs
func vpcListRun() common.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := common.NewCloudConfig(viper.GetViper())
		orgId := config.Vconfig.GetString("org")

		msg := fmt.Sprintf("getting vpcs with env id: %s and org id: %s", env, orgId)
		model := fetch.NewModel(msg, func() (interface{}, int, error) {
			return config.Cc.ListAssets(orgId, env)
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
		dsAssetTypes := []string{"vpc"}
		unfilteredResults := rawResult.Result.([]cloudapiclient.AssetOutput)
		filteredResults := make([]cloudapiclient.AssetOutput, 0)
		for _, result := range unfilteredResults {
			for _, acceptedDsType := range dsAssetTypes {
				if strings.Contains(result.Asset, acceptedDsType) {
					filteredResults = append(filteredResults, result)
				}
			}
		}
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
		Use:     "vpc",
		Short:   "the vpc subcommand helps manage your Aptible vpcs.",
		Long:    `The vpc subcommand helps manage your Aptible vpcs.`,
		Aliases: []string{"v"},
	}

	vpcCreateCmd := &cobra.Command{
		Use:     "create",
		Short:   "provision a new vpc.",
		Long:    `The vpc create command will provision a new vpc.`,
		Aliases: []string{"c", "deploy"},
		RunE:    vpcCreateRun(),
	}

	vpcDestroyCmd := &cobra.Command{
		Use:     "destroy",
		Short:   "permanently remove the vpc.",
		Long:    `The vpc destroy command will permanently remove the vpc.`,
		Aliases: []string{"d", "delete", "rm", "remove"},
		Args:    cobra.MinimumNArgs(1),
		RunE:    vpcDestroyRun(),
	}

	vpcListCmd := &cobra.Command{
		Use:     "list",
		Short:   "list all vpcs within an organization.",
		Long:    `The vpc list command will list all vpcs within an organization.`,
		Aliases: []string{"ls"},
		RunE:    vpcListRun(),
	}

	vpcDestroyCmd.Flags().StringVar(&env, "env", "", "delete vpc within an environment")

	vpcListCmd.Flags().StringVar(&env, "env", "", "list vpc(s) within an environment")

	vpcCmd.AddCommand(vpcCreateCmd)
	vpcCmd.AddCommand(vpcDestroyCmd)
	vpcCmd.AddCommand(vpcListCmd)

	return vpcCmd
}
