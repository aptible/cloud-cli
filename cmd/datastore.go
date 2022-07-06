package cmd

import (
	"fmt"
	"time"

	apiclient "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/aptible/cloud-cli/ui/fetch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func dsCreateRun() RunE {
	return func(cmd *cobra.Command, args []string) error {
		config := NewCloudConfig(viper.GetViper())
		orgID := config.Vconfig.GetString("org")
		fmt.Println(orgID)

		params := apiclient.AssetInput{
			Asset:        "",
			AssetVersion: "",
		}

		/* env, err := config.Cc.CreateAsset(orgID, envID, params)
		if err != nil {
			return err
		} */

		model := fetch.NewModel("creating datastore", func() (interface{}, error) {
			time.Sleep(2 * time.Second)
			return params, nil
		})

		p := tea.NewProgram(model)
		m, err := p.StartReturningModel()
		if err != nil {
			return err
		}

		n := m.(fetch.Model)
		res := n.Result.(apiclient.AssetInput)

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

	var engine string
	var engineVersion string
	var env string

	dsCreateCmd.Flags().StringVarP(&engine, "engine", "e", "", "the datastore engine, e.g. rds/postgres, rds/mysql, etc.")
	dsCreateCmd.Flags().StringVarP(&engineVersion, "version", "v", "", "the engine version, e.g. 14.2")
	dsListCmd.Flags().StringVar(&env, "env", "", "list datastores within an environment")

	datastoreCmd.AddCommand(dsCreateCmd)
	datastoreCmd.AddCommand(dsDestroyCmd)
	datastoreCmd.AddCommand(dsListCmd)

	return datastoreCmd
}
