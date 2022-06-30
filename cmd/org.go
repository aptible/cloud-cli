package cmd

import (
	"fmt"

	apiclient "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func orgCreateRun() RunE {
	return func(cmd *cobra.Command, args []string) error {
		config := NewCloudConfig(viper.GetViper())
		orgID := config.Vconfig.GetString("org")

		output := make(map[string]interface{})
		params := apiclient.OrganizationInput{
			Name:           args[0],
			BaaStatus:      "pending",
			AwsOu:          "idk",
			ContactDetails: output,
		}
		org, err := config.Cc.CreateOrg(orgID, params)
		if err != nil {
			return err
		}

		fmt.Println(fmt.Sprintf("new org: %s\n", org.Name))
		return nil
	}
}

func orgListRun() RunE {
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

func NewOrgCmd() *cobra.Command {
	orgCmd := &cobra.Command{
		Use:     "org",
		Short:   "The org subcommand helps manage your Aptible organizations.",
		Long:    `The org subcommand helps manage your Aptible organizations.`,
		Aliases: []string{"org", "o"},
	}

	orgCreateCmd := &cobra.Command{
		Use:     "create [org name]",
		Short:   "provision a new org.",
		Long:    `The org create command will provision a new org.`,
		Aliases: []string{"c"},
		Args:    cobra.MinimumNArgs(1),
		RunE:    orgCreateRun(),
	}

	orgListCmd := &cobra.Command{
		Use:     "list",
		Short:   "list all orgs you can access.",
		Long:    `list all orgs you can access.`,
		Aliases: []string{"ls"},
		RunE:    orgListRun(),
	}

	orgCmd.AddCommand(orgCreateCmd)
	orgCmd.AddCommand(orgListCmd)

	return orgCmd
}
