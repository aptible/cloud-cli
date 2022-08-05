package cmd

import (
	"fmt"

	cac "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/aptible/cloud-cli/config"
	"github.com/aptible/cloud-cli/lib/env"
	liborg "github.com/aptible/cloud-cli/lib/org"
	"github.com/aptible/cloud-cli/ui/fetch"
	"github.com/aptible/cloud-cli/ui/form"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// envCreateRun - create an environment
func envCreateRun() config.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := config.NewCloudConfig(viper.GetViper())
		org := config.Vconfig.GetString("org")

		formResult := form.FormResult{Org: org}
		err := liborg.OrgForm(config, &formResult)
		if err != nil {
			return err
		}

		desc := ""
		params := cac.EnvironmentInput{
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

		envTable := libenv.EnvTable(result.Result.(*cac.EnvironmentOutput))
		// TODO - print with tea
		fmt.Println("Created Environment(s)")
		fmt.Println(envTable.View())
		return nil
	}
}

// envDestroyRun - destroy an environment
func envDestroyRun() config.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := config.NewCloudConfig(viper.GetViper())
		org := config.Vconfig.GetString("org")
		env := args[0]

		formResult := form.FormResult{Org: org, Env: env}
		err := libenv.EnvForm(config, &formResult)
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
func envListRun() config.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := config.NewCloudConfig(viper.GetViper())
		org := config.Vconfig.GetString("org")

		formResult := form.FormResult{Org: org}
		err := liborg.OrgForm(config, &formResult)
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

		envTable := libenv.EnvTable(result.Result.([]cac.EnvironmentOutput))
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
