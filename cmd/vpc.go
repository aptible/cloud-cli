package cmd

import (
	"fmt"
	"strings"

	apiclient "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/aptible/cloud-cli/ui/fetch"
	"github.com/evertras/bubble-table/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// dataStoreTable - prints out a table of datastores
func vpcTable(orgOutput interface{}) table.Model {
	rows := make([]table.Row, 0)

	switch data := orgOutput.(type) {
	case []apiclient.AssetOutput:
		for _, asset := range data {
			rows = append(rows, table.NewRow(table.RowData{
				"id":     asset.Id,
				"status": asset.Status,
			}))
		}
	case *apiclient.AssetOutput:
		rows = append(rows, table.NewRow(table.RowData{
			"id":     data.Id,
			"status": data.Status,
		}))
	}

	return table.New([]table.Column{
		table.NewColumn("id", "VPC Asset Id", 40),
		table.NewColumn("status", "VPC Status", 40),
	}).WithRows(rows)
}

// dsCreateRun - create a datastore
func vpcCreateRun() CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := NewCloudConfig(viper.GetViper())
		orgId := config.Vconfig.GetString("org")
		envId := args[0]
		name := args[1]

		vars := map[string]interface{}{
			"name": name,
		}
		params := apiclient.AssetInput{
			Asset:           "aws__vpc__latest",
			AssetVersion:    "latest",
			AssetParameters: vars,
		}

		msg := fmt.Sprintf("creating vpc %s (v%s)", engine, engineVersion)
		model := fetch.NewModel(msg, func() (interface{}, error) {
			return config.Cc.CreateAsset(orgId, envId, params)
		})

		result, err := fetch.FetchWithOutput(model)
		if err != nil {
			return err
		}
		vpcTable := vpcTable(result.Result.(*apiclient.AssetOutput))
		// TODO - print with tea
		fmt.Println("VPC(s) List")
		fmt.Println(vpcTable.View())

		return nil
	}
}

// dsDestroyRun - destroy datastore
func vpcDestroyRun() CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		fmt.Println(fmt.Sprintf("Destroying vpc id: %s", args[0]))
		return destroyAsset(cmd, args)
	}
}

// vpcListRun - list vpcs
func vpcListRun() CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := NewCloudConfig(viper.GetViper())
		orgId := config.Vconfig.GetString("org")

		msg := fmt.Sprintf("getting vpcs with env id: %s and org id: %s", env, orgId)
		model := fetch.NewModel(msg, func() (interface{}, error) {
			return config.Cc.ListAssets(orgId, env)
		})

		rawResult, err := fetch.FetchWithOutput(model)
		if err != nil {
			return err
		}
		if rawResult == nil {
			fmt.Println("No datastores found.")
			return nil
		}
		dsAssetTypes := []string{"vpc"}
		unfilteredResults := rawResult.Result.([]apiclient.AssetOutput)
		filteredResults := make([]apiclient.AssetOutput, 0)
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
		Aliases: []string{"database", "v"},
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
		Short:   "list all datastores within an organization.",
		Long:    `The datastore list command will list all datastores within an organization.`,
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
