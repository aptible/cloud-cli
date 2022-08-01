package cmd

import (
	"fmt"

	cloudapiclient "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aptible/cloud-cli/internal/common"
	"github.com/aptible/cloud-cli/internal/ui/fetch"
	"github.com/aptible/cloud-cli/internal/ui/form"
	"github.com/aptible/cloud-cli/table"
)

// envCreateRun - create an environment
func envCreateRun() common.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := common.NewCloudConfig(viper.GetViper())
		org := config.Vconfig.GetString("org")

		formResult := form.FormResult{Org: org}
		err := form.OrgForm(config, &formResult)
		if err != nil {
			return err
		}

		desc := ""
		params := cloudapiclient.EnvironmentInput{
			Name:        args[0],
			Description: &desc,
			Data:        map[string]interface{}{},
		}

		progressModel := fetch.NewModel("creating environment", func() (interface{}, error) {
			return config.Cc.CreateEnvironment(formResult.Org, params)
		})

		result, err := fetch.WithOutput(progressModel)
		if err != nil {
			return err
		}

		envTable := table.EnvTable(result.Result.(*cloudapiclient.EnvironmentOutput))
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
		org := config.Vconfig.GetString("org")
		env := args[0]

		formResult := form.FormResult{Org: org, Env: env}
		err := form.EnvForm(config, &formResult)
		if err != nil {
			return err
		}

		model := fetch.NewModel("destroying environment", func() (interface{}, error) {
			err := config.Cc.DestroyEnvironment(formResult.Org, formResult.Env)
			return nil, err
		})

		err = fetch.Any(model)

		// does not print anything, no table to print here
		fmt.Printf("Destroyed environment: %s\n", env)
		return err
	}
}

// envListRun - lists all environments for an org id
func envListRun() common.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := common.NewCloudConfig(viper.GetViper())
		org := config.Vconfig.GetString("org")

		formResult := form.FormResult{Org: org}
		err := form.OrgForm(config, &formResult)
		if err != nil {
			return err
		}

		model := fetch.NewModel("fetching environments", func() (interface{}, error) {
			return config.Cc.ListEnvironments(formResult.Org)
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

		envTable := table.EnvTable(result.Result.([]cloudapiclient.EnvironmentOutput))
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
		Use:     "create [env_name]",
		Short:   "creates an environment under the current organization.",
		Long:    `The environment create command will provision a new environment.`,
		Aliases: []string{"c"},
		Args:    cobra.MinimumNArgs(1),
		RunE:    envCreateRun(),
	}

	envDestroyCmd := &cobra.Command{
		Use:     "destroy [env_id]",
		Short:   "permentantly remove the environment.",
		Long:    `The environment destroy command will permanently remove the environment.`,
		Aliases: []string{"d", "delete", "rm", "remove"},
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
