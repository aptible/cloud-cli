package cmd

import (
	"fmt"

	apiclient "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/aptible/cloud-cli/ui/fetch"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	engine        string
	engineVersion string
	env           string
)

func dsCreateRun() RunE {
	return func(cmd *cobra.Command, args []string) error {
		config := NewCloudConfig(viper.GetViper())
		orgID := config.Vconfig.GetString("org")
		envID := args[0]
		name := args[1]

		if engine == "" {
			return fmt.Errorf("must provide engine")
		}
		if engineVersion == "" {
			return fmt.Errorf("must provide engine version")
		}

		vars := map[string]interface{}{
			"name":           name,
			"engine":         engine,
			"engine_version": engineVersion,
		}
		params := apiclient.AssetInput{
			Asset:           "aws__rds__latest",
			AssetVersion:    "latest",
			AssetParameters: vars,
		}

		msg := fmt.Sprintf("creating datastore %s (v%s)", engine, engineVersion)
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

func dsDestroyRun() RunE {
	return func(cmd *cobra.Command, args []string) error {
		config := NewCloudConfig(viper.GetViper())
		orgID := config.Vconfig.GetString("org")
		fmt.Println(orgID)
		return nil
	}
}

func dsListRun() RunE {
	return func(cmd *cobra.Command, args []string) error {
		config := NewCloudConfig(viper.GetViper())
		orgID := config.Vconfig.GetString("org")
		fmt.Println(orgID)
		return nil
	}
}

func NewDatastoreCmd() *cobra.Command {
	datastoreCmd := &cobra.Command{
		Use:     "datastore",
		Short:   "The datastore subcommand helps manage your Aptible datastores.",
		Long:    `The datastore subcommand helps manage your Aptible datastores.`,
		Aliases: []string{"database", "ds", "db"},
	}

	dsCreateCmd := &cobra.Command{
		Use:     "create",
		Short:   "provision a new datastore.",
		Long:    `The datastore create command will provision a new datastore.`,
		Aliases: []string{"c", "deploy"},
		Args:    cobra.MinimumNArgs(2),
		RunE:    dsCreateRun(),
	}

	dsDestroyCmd := &cobra.Command{
		Use:     "destroy",
		Short:   "permentantly remove the datastore.",
		Long:    `The datastore destroy command will permentantly remove the datastore.`,
		Aliases: []string{"d", "delete", "rm", "remove"},
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
	dsCreateCmd.Flags().StringVarP(&engineVersion, "version", "v", "", "the engine version, e.g. 14.2")
	dsListCmd.Flags().StringVar(&env, "env", "", "list datastores within an environment")

	datastoreCmd.AddCommand(dsCreateCmd)
	datastoreCmd.AddCommand(dsDestroyCmd)
	datastoreCmd.AddCommand(dsListCmd)

	return datastoreCmd
}
