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
	"github.com/aptible/cloud-cli/internal/ui/form"
)

// organizationsTable - prints out a table of organizations
func organizationsTable(orgOutput interface{}) table.Model {
	rows := make([]table.Row, 0)

	switch data := orgOutput.(type) {
	case []cloudapiclient.OrganizationOutput:
		for _, org := range data {
			rows = append(rows, table.NewRow(table.RowData{
				"id":     org.Id,
				"name":   org.Name,
				"aws_ou": *org.AwsOu,
			}))
		}
	case *cloudapiclient.OrganizationOutput:
		rows = append(rows, table.NewRow(table.RowData{
			"id":     data.Id,
			"name":   data.Name,
			"aws_ou": *data.AwsOu,
		}))
	}

	return table.New([]table.Column{
		table.NewColumn("id", "Organization Id", 40).WithStyle(uiCommon.DefaultRowStyle()),
		table.NewColumn("name", "Organization Name", 40).WithStyle(uiCommon.DefaultRowStyle()),
		table.NewColumn("aws_ou", "AWS OU", 40).WithStyle(uiCommon.DefaultRowStyle()),
	}).WithRows(rows)
}

// organizationCreateRun - create an organization
func organizationCreateRun() common.CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := common.NewCloudConfig(viper.GetViper())
		orgId := config.Vconfig.GetString("org")

		formResult, err := form.OrgForm(config, orgId)
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

		orgsTable := organizationsTable(result.Result.(*cloudapiclient.OrganizationOutput))
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

		orgsTable := organizationsTable(result.Result.([]cloudapiclient.OrganizationOutput))
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
