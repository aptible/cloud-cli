package asset

import (
	"fmt"
	"strings"

	cac "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/aptible/cloud-cli/config"
	libasset "github.com/aptible/cloud-cli/lib/asset"
	libenv "github.com/aptible/cloud-cli/lib/env"
	"github.com/aptible/cloud-cli/ui/fetch"
	"github.com/aptible/cloud-cli/ui/form"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	engine        string
	engineVersion string
	name          string
	vpcName       string
)

// dsCreateRun - create a datastore
func dsCreateRun() config.CobraRunE {
	return assetsCreateRun()
}

// dsDescribeRun - describe datastore
func dsDescribeRun() config.CobraRunE {
	return describeAsset()
}

// dsDestroyRun - destroy datastore
func dsDestroyRun() config.CobraRunE {
	return destroyAsset()
}

// dsListRun - list datastores
func dsListRun() config.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := config.NewCloudConfig(viper.GetViper())
		org := config.Vconfig.GetString("org")
		env := config.Vconfig.GetString("env")

		formResult := form.FormResult{Org: org, Env: env}
		err := libenv.EnvForm(config, &formResult)
		if err != nil {
			return nil
		}

		msg := fmt.Sprintf("geting datastores with %+v", formResult)
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

		dsAssetTypes := []string{"rds"}
		unfilteredResults := rawResult.Result.([]cac.AssetOutput)
		filteredResults := make([]cac.AssetOutput, 0)
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

		dsTable := libasset.AssetTable(filteredResults)
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
		Use:     "destroy [datastore_id]",
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
		RunE:    dsDescribeRun(),
	}

	dsCreateCmd.Flags().StringVarP(&engine, "engine", "e", "", "the datastore engine, e.g. postgres, mysql, etc.")
	dsCreateCmd.Flags().StringVarP(&engineVersion, "engine-version", "v", "", "the engine version, e.g. 14.2")
	dsCreateCmd.Flags().StringVar(&name, "name", "", "the name to assign to rds")
	dsCreateCmd.Flags().StringVarP(&vpcName, "vpc-name", "", "", "the vpc to attach rds to")

	datastoreCmd.AddCommand(dsCreateCmd)
	datastoreCmd.AddCommand(dsDestroyCmd)
	datastoreCmd.AddCommand(dsListCmd)
	datastoreCmd.AddCommand(dsDescribeCmd)

	return datastoreCmd
}
