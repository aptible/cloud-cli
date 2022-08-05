package cmd

import (
	"fmt"

	cac "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/aptible/cloud-cli/config"
	"github.com/aptible/cloud-cli/lib/org"
	"github.com/aptible/cloud-cli/ui/fetch"
	"github.com/aptible/cloud-cli/ui/form"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// organizationCreateRun - create an organization
func organizationCreateRun() config.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := config.NewCloudConfig(viper.GetViper())
		org := config.Vconfig.GetString("org")

		formResult := form.FormResult{Org: org, Env: env}
		err := liborg.OrgForm(config, &formResult)
		if err != nil {
			return err
		}

		output := make(map[string]interface{})
		var ou string
		params := cac.OrganizationInput{
			Name:           args[0],
			BaaStatus:      "pending",
			AwsOu:          &ou,
			ContactDetails: output,
		}

		progressModel := fetch.NewModel("creating organization", func() (interface{}, error) {
			return config.Cc.CreateOrg(formResult.Org, params)
		})
		result, err := fetch.WithOutput(progressModel)
		if err != nil {
			return err
		}

		orgsTable := liborg.OrgTable(result.Result.(*cac.OrganizationOutput))
		// TODO - print with tea
		fmt.Println("Created Organization(s)")
		fmt.Println(orgsTable.View())

		return nil
	}
}

// orgListRun - lists all organizations
func orgListRun() config.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := config.NewCloudConfig(viper.GetViper())
		progressModel := fetch.NewModel("fetching organizations", func() (interface{}, error) {
			return config.Cc.ListOrgs()
		})
		result, err := fetch.WithOutput(progressModel)
		if err != nil {
			return err
		}
		if result == nil {
			fmt.Println("No organizations found.")
			return nil
		}

		orgsTable := liborg.OrgTable(result.Result.([]cac.OrganizationOutput))
		// TODO - print with tea
		fmt.Println("Organization(s) List")
		fmt.Println(orgsTable.View())

		return nil
	}
}

func NewOrgCmd() *cobra.Command {
	orgCmd := &cobra.Command{
		Use:     "organization",
		Short:   "The org subcommand helps manage your Aptible organizations.",
		Long:    `The org subcommand helps manage your Aptible organizations.`,
		Aliases: []string{"org", "o"},
	}

	orgCreateCmd := &cobra.Command{
		Use:     "create [org_name]",
		Short:   "provision a new org.",
		Long:    `The org create command will provision a new org.`,
		Aliases: []string{"c"},
		Args:    cobra.MinimumNArgs(1),
		RunE:    organizationCreateRun(),
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
