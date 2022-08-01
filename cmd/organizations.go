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

// organizationCreateRun - create an organization
func organizationCreateRun() common.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := common.NewCloudConfig(viper.GetViper())
		org := config.Vconfig.GetString("org")

		formResult := form.FormResult{Org: org, Env: env}
		err := form.OrgForm(config, &formResult)
		if err != nil {
			return err
		}

		output := make(map[string]interface{})
		var ou string
		params := cloudapiclient.OrganizationInput{
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

		orgsTable := table.OrgTable(result.Result.(*cloudapiclient.OrganizationOutput))
		// TODO - print with tea
		fmt.Println("Created Organization(s)")
		fmt.Println(orgsTable.View())

		return nil
	}
}

// orgListRun - lists all organizations
func orgListRun() common.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := common.NewCloudConfig(viper.GetViper())
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

		orgsTable := table.OrgTable(result.Result.([]cloudapiclient.OrganizationOutput))
		// TODO - print with tea
		fmt.Println("Organization(s) List")
		fmt.Println(orgsTable.View())

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
