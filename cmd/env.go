package cmd

import (
	"fmt"

	apiclient "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/aptible/cloud-cli/ui/fetch"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func envCreateRun() RunE {
	return func(cmd *cobra.Command, args []string) error {
		config := NewCloudConfig(viper.GetViper())
		orgID := config.Vconfig.GetString("org")
		desc := ""
		params := apiclient.EnvironmentInput{
			Name:        args[0],
			Description: &desc,
			Data:        map[string]interface{}{},
		}

		model := fetch.NewModel("creating environment", func() (interface{}, error) {
			return config.Cc.CreateEnvironment(orgID, params)
		})

		fetchModel, err := fetch.FetchWithOutput(model)
		if err != nil {
			return err
		}
		env := fetchModel.Result.(apiclient.EnvironmentOutput)

		fmt.Printf("New environment ID: %s\n", env.Id)
		return nil
	}
}

func envDestroyRun() RunE {
	return func(cmd *cobra.Command, args []string) error {
		config := NewCloudConfig(viper.GetViper())
		orgID := config.Vconfig.GetString("org")
		envID := args[0]

		model := fetch.NewModel("destroying environment", func() (interface{}, error) {
			err := config.Cc.DestroyEnvironment(orgID, envID)
			return nil, err
		})

		err := fetch.FetchAny(model)
		return err
	}
}

func envListRun() RunE {
	return func(cmd *cobra.Command, args []string) error {
		config := NewCloudConfig(viper.GetViper())
		orgID := config.Vconfig.GetString("org")

		model := fetch.NewModel("fetching environments", func() (interface{}, error) {
			return config.Cc.ListEnvironments(orgID)
		})
		result, err := fetch.FetchWithOutput(model)
		if err != nil {
			return err
		}

		envs := result.Result.([]apiclient.EnvironmentOutput)

		for _, env := range envs {
			fmt.Printf("%s %s\n", env.Id, env.Name)
		}
		return nil
	}
}

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
		Long:    `The datastore destroy command will permentantly remove the environment.`,
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
