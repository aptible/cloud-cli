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

var (
	engine        string
	engineVersion string
	name          string
	vpcName       string
)

func generateDatastoreRowFromData(asset cloudapiclient.AssetOutput) table.Row {
	row := table.NewRow(table.RowData{
		"id":             asset.Id,
		"status":         asset.Status,
		"name":           asset.CurrentAssetParameters.Data["name"],
		"engine":         asset.CurrentAssetParameters.Data["engine"],
		"engine_version": asset.CurrentAssetParameters.Data["engine_version"],
		"vpc_name":       asset.CurrentAssetParameters.Data["vpc_name"],
	})
	return colorizeAssetFromStatus(asset, row)
}

// dataStoreTable - prints out a table of datastores
func dataStoreTable(orgOutput interface{}) table.Model {
	rows := make([]table.Row, 0)

	switch data := orgOutput.(type) {
	case []cloudapiclient.AssetOutput:
		for _, asset := range data {
			rows = append(rows, generateDatastoreRowFromData(asset))
		}
	case *cloudapiclient.AssetOutput:
		rows = append(rows, generateDatastoreRowFromData(*data))
	}

	return table.New([]table.Column{
		table.NewColumn("id", "Id", 40).WithStyle(uiCommon.DefaultRowStyle()),
		table.NewColumn("status", "Status", 40).WithStyle(uiCommon.DefaultRowStyle()),
		table.NewColumn("name", "Name", 20).WithStyle(uiCommon.DefaultRowStyle()),
		table.NewColumn("engine", "Engine", 20).WithStyle(uiCommon.DefaultRowStyle()),
		table.NewColumn("engine_version", "Engine Version", 20).WithStyle(uiCommon.DefaultRowStyle()),
		table.NewColumn("vpc_name", "VPC Name", 20).WithStyle(uiCommon.DefaultRowStyle()),
	}).WithRows(rows)
}

// dsCreateRun - create a datastore
func dsCreateRun() common.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := common.NewCloudConfig(viper.GetViper())
		orgId := config.Vconfig.GetString("org")
		envId := config.Vconfig.GetString("env")

		if envId == "" {
			return fmt.Errorf("must provide env")
		}
		if engine == "" {
			return fmt.Errorf("must provide engine")
		}
		if engineVersion == "" {
			return fmt.Errorf("must provide engine version")
		}
		if name == "" {
			return fmt.Errorf("must provide name")
		}
		if vpcName == "" {
			return fmt.Errorf("must provide vpc-name")
		}

		vars := map[string]interface{}{
			"name":           name,
			"engine":         engine,
			"engine_version": engineVersion,
			"vpc_name":       vpcName,
		}
		params := cloudapiclient.AssetInput{
			Asset:           "aws__rds__latest",
			AssetVersion:    "latest",
			AssetParameters: vars,
		}

		msg := fmt.Sprintf("creating datastore %s (v%s)", engine, engineVersion)
		model := fetch.NewModel(msg, func() (interface{}, int, error) {
			return config.Cc.CreateAsset(orgId, envId, params)
		})

		result, err := fetch.WithOutput(model)
		if err != nil {
			return err
		}
		datastoreTable := vpcTable(result.Result.(*cloudapiclient.AssetOutput))
		// TODO - print with tea
		fmt.Println("Datastore(s) Created:")
		fmt.Println(datastoreTable.View())

		return nil
	}
}

// dsDescribeRun - describe datastore
func dsDescribeRun() common.CobraRunE {
	return describeAsset()
}

// dsDestroyRun - destroy datastore
func dsDestroyRun() common.CobraRunE {
	return destroyAsset()
}

// dsListRun - list datastores
func dsListRun() common.CobraRunE {
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

		dsAssetTypes := []string{"rds"}
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
			// TODO - print with tea
			fmt.Println("No datastores found.")
			return nil
		}

		dsTable := dataStoreTable(filteredResults)
		// TODO - print with tea
		fmt.Println("Datastore(s) List")
		fmt.Println(dsTable.View())

		return nil
	}
}

func NewDatastoreCmd() *cobra.Command {
	datastoreCmd := &cobra.Command{
		Use:     "datastore",
		Short:   "the datastore subcommand helps manage your Aptible datastore assets.",
		Long:    `The datastore subcommand helps manage your Aptible datastore assets.`,
		Aliases: []string{"database", "ds", "db", "rds"},
	}

	dsCreateCmd := &cobra.Command{
		Use:     "create",
		Short:   "provision a new datastore.",
		Long:    `The datastore create command will provision a new datastore.`,
		Aliases: []string{"c", "deploy"},
		RunE:    dsCreateRun(),
	}

	dsDestroyCmd := &cobra.Command{
		Use:     "destroy",
		Short:   "permanently remove the datastore.",
		Long:    `The datastore destroy command will permanently remove the datastore.`,
		Aliases: []string{"d", "delete", "rm", "remove"},
		Args:    cobra.MinimumNArgs(1),
		RunE:    dsDestroyRun(),
	}

	dsListCmd := &cobra.Command{
		Use:     "list",
		Short:   "list all datastores within an organization.",
		Long:    `The datastore list command will list all datastores within an organization.`,
		Aliases: []string{"ls"},
		RunE:    dsListRun(),
	}

	dsDescribeCmd := &cobra.Command{
		Use:     "describe",
		Short:   "describe datastore",
		Long:    `The datastore show command will provide more detail about a datastore`,
		Aliases: []string{"show"},
		RunE:    dsListRun(),
	}

	dsCreateCmd.Flags().StringVarP(&engine, "engine", "e", "", "the datastore engine, e.g. rds/postgres, rds/mysql, etc.")
	dsCreateCmd.Flags().StringVarP(&engineVersion, "engine-version", "v", "", "the engine version, e.g. 14.2")
	dsCreateCmd.Flags().StringVar(&name, "name", "", "the name to assign to rds")
	dsCreateCmd.Flags().StringVarP(&vpcName, "vpc-name", "", "", "the vpc to attach rds to")

	datastoreCmd.AddCommand(dsCreateCmd)
	datastoreCmd.AddCommand(dsDestroyCmd)
	datastoreCmd.AddCommand(dsListCmd)
	datastoreCmd.AddCommand(dsDescribeCmd)

	return datastoreCmd
}
