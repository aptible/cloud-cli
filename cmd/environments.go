package cmd

import (
	"fmt"

	cloudapiclient "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/evertras/bubble-table/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aptible/cloud-cli/internal/common"
	uiCommon "github.com/aptible/cloud-cli/internal/ui/common"
	"github.com/aptible/cloud-cli/internal/ui/fetch"
)

// environmentsTable - prints out a table of environments
func environmentsTable(orgOutput interface{}) table.Model {
	rows := make([]table.Row, 0)

	switch data := orgOutput.(type) {
	case []cloudapiclient.EnvironmentOutput:
		for _, org := range data {
			rows = append(rows, table.NewRow(table.RowData{
				"id":   org.Id,
				"name": org.Name,
			}))
		}
	case cloudapiclient.EnvironmentOutput:
		rows = append(rows, table.NewRow(table.RowData{
			"id":   data.Id,
			"name": data.Name,
		}))
	}

	return table.New([]table.Column{
		table.NewColumn("id", "Environment Id", 40).WithStyle(uiCommon.DefaultRowStyle()),
		table.NewColumn("name", "Environment Name", 40).WithStyle(uiCommon.DefaultRowStyle()),
	}).WithRows(rows)
}

// envCreateRun - create an environment
func envCreateRun() common.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := common.NewCloudConfig(viper.GetViper())
		orgId := config.Vconfig.GetString("org")
		desc := ""
		params := cloudapiclient.EnvironmentInput{
			Name:        args[0],
			Description: &desc,
			Data:        map[string]interface{}{},
		}

		progressModel := fetch.NewModel("creating environment", func() (interface{}, int, error) {
			return config.Cc.CreateEnvironment(orgId, params)
		})

		result, err := fetch.WithOutput(progressModel)
		if err != nil {
			return err
		}

		envTable := environmentsTable(result.Result.(cloudapiclient.EnvironmentOutput))
		// TODO - print with tea
		fmt.Println("Created Environment(s)")
		fmt.Println(envTable.View())
		return nil
	}
}

// envDestroyRun - destroy an environment
func envDestroyRun() common.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := common.NewCloudConfig(viper.GetViper())
		orgId := config.Vconfig.GetString("org")
		envId := args[0]

		model := fetch.NewModel("destroying environment", func() (interface{}, int, error) {
			status, err := config.Cc.DestroyEnvironment(orgId, envId)
			return nil, status, err
		})

		err := fetch.Any(model)

		// does not print anything, no table to print here
		fmt.Println(fmt.Sprintf("Destroyed environment: %s", envId))
		return err
	}
}

// envListRun - lists all environments for an org id
func envListRun() common.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := common.NewCloudConfig(viper.GetViper())
		orgId := config.Vconfig.GetString("org")
		model := fetch.NewModel("fetching environments", func() (interface{}, int, error) {
			return config.Cc.ListEnvironments(orgId)
		})
		result, err := fetch.WithOutput(model)
		if err != nil {
			return err
		}
		if result == nil {
			// TODO - print with tea
			fmt.Println("No environments found.")
			return nil
		}

		envTable := environmentsTable(result.Result.([]cloudapiclient.EnvironmentOutput))
		// TODO - print with tea
		fmt.Println("Environment(s) List")
		fmt.Println(envTable.View())

		return nil
	}
}

// NewEnvCmd - generates a cobra command target for environments
func NewEnvCmd() *cobra.Command {
	envCmd := &cobra.Command{
		Use:     "environment",
		Short:   "The env subcommand helps manage your Aptible environments.",
		Long:    `The env subcommand helps manage your Aptible environments.`,
		Aliases: []string{"env", "e"},
	}

	envCreateCmd := &cobra.Command{
		Use:     "create",
		Short:   "provision a new datastore.",
		Long:    `The environment create command will provision a new environment.`,
		Aliases: []string{"c"},
		Args:    cobra.MinimumNArgs(1),
		RunE:    envCreateRun(),
	}

	envDestroyCmd := &cobra.Command{
		Use:     "destroy",
		Short:   "permentantly remove the environment.",
		Long:    `The datastore destroy command will permanently remove the environment.`,
		Aliases: []string{"d", "delete", "rm", "remove"},
		Args:    cobra.MinimumNArgs(1),
		RunE:    envDestroyRun(),
	}

	envListCmd := &cobra.Command{
		Use:     "list",
		Short:   "list all environment within an organization.",
		Long:    `The environment list command will list all environment within an organization.`,
		Aliases: []string{"ls"},
		RunE:    envListRun(),
	}

	envCmd.AddCommand(envCreateCmd)
	envCmd.AddCommand(envDestroyCmd)
	envCmd.AddCommand(envListCmd)

	return envCmd
}
