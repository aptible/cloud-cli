package cmd

import (
	"fmt"
	"time"

	apiclient "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/aptible/cloud-cli/ui/fetch"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	tea "github.com/charmbracelet/bubbletea"
)

func envCreateRun() RunE {
	return func(cmd *cobra.Command, args []string) error {
		// config := NewCloudConfig(viper.GetViper())
		// orgID := config.Vconfig.GetString("org")
		params := apiclient.EnvironmentInput{
			Name: args[0],
		}

		/* env, err := config.Cc.CreateEnvironment(orgID, params)
		if err != nil {
			return err
		} */

		fx := func() (interface{}, error) {
			time.Sleep(2 * time.Second)
			return params, nil
		}
		model := fetch.NewModel(fx, "creating environment")

		p := tea.NewProgram(model)
		m, err := p.StartReturningModel()
		if err != nil {
			return err
		}

		n := m.(fetch.Model)
		res := n.Result.(apiclient.EnvironmentInput)

		fmt.Printf("Result: %+v\n", res)
		return nil
	}
}

func envDestroyRun() RunE {
	return func(cmd *cobra.Command, args []string) error {
		config := NewCloudConfig(viper.GetViper())
		orgID := config.Vconfig.GetString("org")
		envID := ""
		err := config.Cc.DestroyEnvironment(orgID, envID)
		if err != nil {
			return err
		}
		return nil
	}
}

func envListRun() RunE {
	return func(cmd *cobra.Command, args []string) error {
		config := NewCloudConfig(viper.GetViper())
		orgID := config.Vconfig.GetString("org")
		envs, err := config.Cc.ListEnvironments(orgID)
		if err != nil {
			return err
		}

		for _, env := range envs {
			fmt.Println(fmt.Println(env.Name))
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
