package cmd

import (
	"fmt"

	apiclient "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func envCreateRun(vconfig *viper.Viper) RunE {
	return func(cmd *cobra.Command, args []string) error {
		config := NewCloudConfig(vconfig)
		orgID := vconfig.GetString("org")
		params := apiclient.EnvironmentInput{
			Name: args[0],
		}
		env, err := config.Cc.CreateEnvironment(orgID, params)
		if err != nil {
			return err
		}

		fmt.Println(fmt.Sprintf("new env id: %s\n", env.Id))
		return nil
	}
}

func envDestroyRun(vconfig *viper.Viper) RunE {
	return func(cmd *cobra.Command, args []string) error {
		config := NewCloudConfig(vconfig)
		orgID := vconfig.GetString("org")
		envID := ""
		err := config.Cc.DestroyEnvironment(orgID, envID)
		if err != nil {
			return err
		}
		return nil
	}
}

func envListRun(vconfig *viper.Viper) RunE {
	return func(cmd *cobra.Command, args []string) error {
		config := NewCloudConfig(vconfig)
		orgID := vconfig.GetString("org")
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

func NewEnvCmd(vconfig *viper.Viper) *cobra.Command {
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
		RunE:    envCreateRun(vconfig),
	}

	envDestroyCmd := &cobra.Command{
		Use:     "destroy",
		Short:   "permentantly remove the environment.",
		Long:    `The datastore destroy command will permentantly remove the environment.`,
		Aliases: []string{"d", "delete", "rm", "remove"},
		RunE:    envDestroyRun(vconfig),
	}

	envListCmd := &cobra.Command{
		Use:     "list",
		Short:   "list all environment within an organization.",
		Long:    `The environment list command will list all environment within an organization.`,
		Aliases: []string{"ls"},
		RunE:    envListRun(vconfig),
	}

	envCmd.AddCommand(envCreateCmd)
	envCmd.AddCommand(envDestroyCmd)
	envCmd.AddCommand(envListCmd)

	return envCmd
}
