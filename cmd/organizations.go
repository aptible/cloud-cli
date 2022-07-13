package cmd

import (
	"fmt"

	apiclient "github.com/aptible/cloud-api-clients/clients/go"
	"github.com/aptible/cloud-cli/ui/fetch"
	"github.com/evertras/bubble-table/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// organizationsTable - prints out a table of organizations
func organizationsTable(orgOutput interface{}) table.Model {
	rows := make([]table.Row, 0)

	switch data := orgOutput.(type) {
	case []apiclient.OrganizationOutput:
		for _, org := range data {
			rows = append(rows, table.NewRow(table.RowData{
				"id":   org.Id,
				"name": org.Name,
			}))
		}
	case apiclient.OrganizationOutput:
		rows = append(rows, table.NewRow(table.RowData{
			"id":   data.Id,
			"name": data.Name,
		}))
	}

	return table.New([]table.Column{
		table.NewColumn("id", "Organization Id", 40),
		table.NewColumn("name", "Organization Name", 40),
	}).WithRows(rows)
}

// organizationCreateRun - create an organization
func organizationCreateRun() CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := NewCloudConfig(viper.GetViper())
		orgId := config.Vconfig.GetString("org")

		output := make(map[string]interface{})
		var ou string
		params := apiclient.OrganizationInput{
			Name:           args[0],
			BaaStatus:      "pending",
			AwsOu:          &ou,
			ContactDetails: output,
		}

		progressModel := fetch.NewModel("creating organization", func() (interface{}, error) {
			return config.Cc.CreateOrg(orgId, params)
		})
		result, err := fetch.FetchWithOutput(progressModel)
		if err != nil {
			return err
		}

		orgsTable := organizationsTable(result.Result.(apiclient.OrganizationOutput))
		// TODO - print with tea
		fmt.Println("Created Organization(s)")
		fmt.Println(orgsTable.View())

		return nil
	}
}

// orgListRun - lists all organizations
func orgListRun() CobraRunE {
	return func(cmd *cobra.Command, args []string) error {
		config := NewCloudConfig(viper.GetViper())
		progressModel := fetch.NewModel("fetching organizations", func() (interface{}, error) {
			return config.Cc.ListOrgs()
		})
		result, err := fetch.FetchWithOutput(progressModel)
		if err != nil {
			return err
		}
		if result == nil {
			fmt.Println("No organizations found.")
			return nil
		}

		orgsTable := organizationsTable(result.Result.([]apiclient.OrganizationOutput))
		// TODO - print with tea
		fmt.Println("Organizations List")
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
		Use:     "create [org name]",
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
