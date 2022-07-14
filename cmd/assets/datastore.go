package assets

import (
	"fmt"
	"github.com/aptible/cloud-cli/internal/common"
	"strings"

	cloudapiclient "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/evertras/bubble-table/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aptible/cloud-cli/internal/ui/fetch"
)

var (
	engine        string
	engineVersion string
	env           string
	name          string
	vpcName       string
)

// dataStoreTable - prints out a table of datastores
func dataStoreTable(orgOutput interface{}) table.Model {
	rows := make([]table.Row, 0)

	switch data := orgOutput.(type) {
	case []cloudapiclient.AssetOutput:
		for _, asset := range data {
			rows = append(rows, table.NewRow(table.RowData{
				"id":     asset.Id,
				"status": asset.Status,
			}))
		}
	case cloudapiclient.AssetOutput:
		rows = append(rows, table.NewRow(table.RowData{
			"id":     data.Id,
			"status": data.Status,
		}))
	}

	return table.New([]table.Column{
		table.NewColumn("id", "Datastore Id", 40),
		table.NewColumn("status", "Datastore Status", 40),
	}).WithRows(rows)
}

// dsCreateRun - create a datastore
func dsCreateRun() common.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := common.NewCloudConfig(viper.GetViper())
		orgId := config.Vconfig.GetString("org")

		if env == "" {
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
			return config.Cc.CreateAsset(orgId, env, params)
		})

		result, err := fetch.FetchWithOutput(model)
		if err != nil {
			return err
		}
		res := result.Result.(*cloudapiclient.AssetOutput)

		fmt.Printf("Result: %+v\n", res)
		return nil
	}
}

// dsDestroyRun - destroy datastore
func dsDestroyRun() common.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		fmt.Println(fmt.Sprintf("Destroying datastore id: %s", args[0]))
		return destroyAsset(cmd, args)
	}
}

// dsListRun - list datastores
func dsListRun() common.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := common.NewCloudConfig(viper.GetViper())
		orgId := config.Vconfig.GetString("org")

		msg := fmt.Sprintf("getting datastores with env id: %s and org id: %s", env, orgId)
		model := fetch.NewModel(msg, func() (interface{}, int, error) {
			return config.Cc.ListAssets(orgId, env)
		})

		rawResult, err := fetch.FetchWithOutput(model)
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
		Short:   "the datastore subcommand helps manage your Aptible datastores.",
		Long:    `The datastore subcommand helps manage your Aptible datastores.`,
		Aliases: []string{"database", "ds", "db"},
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

	dsCreateCmd.Flags().StringVarP(&engine, "engine", "e", "", "the datastore engine, e.g. rds/postgres, rds/mysql, etc.")
	dsCreateCmd.Flags().StringVarP(&engineVersion, "engine-version", "v", "", "the engine version, e.g. 14.2")
	dsCreateCmd.Flags().StringVar(&env, "env", "", "the environment id to deploy to")
	dsCreateCmd.Flags().StringVar(&name, "name", "", "the name to assign to rds")
	dsCreateCmd.Flags().StringVarP(&vpcName, "vpc-name", "", "", "the vpc to attach rds to")

	dsDestroyCmd.Flags().StringVar(&env, "env", "", "delete datastore within an environment")

	dsListCmd.Flags().StringVar(&env, "env", "", "list datastores within an environment")

	datastoreCmd.AddCommand(dsCreateCmd)
	datastoreCmd.AddCommand(dsDestroyCmd)
	datastoreCmd.AddCommand(dsListCmd)

	return datastoreCmd
}
