package form

import (
	"fmt"

	"github.com/aptible/cloud-cli/internal/common"
	"github.com/charmbracelet/bubbles/list"
)

func NewOrgProp() *SubSchema {
	return &SubSchema{
		Type:        "select",
		Title:       "Select an organization",
		LoadOptions: CreateOrgOptions(),
	}
}

func CreateOrgOptions() LoadOptionsFn {
	var options []list.Item
	return func(config *common.CloudConfig) ([]list.Item, error) {
		orgs, err := config.Cc.ListOrgs()
		if err != nil {
			return options, err
		}
		for _, org := range orgs {
			options = append(options, FormOption{Label: org.Name, Value: org.Id})
		}
		return options, nil
	}
}

func OrgForm(config *common.CloudConfig, orgId string) (FormResult, error) {
	results := FormResult{
		Org: orgId,
	}

	if results.Org == "" {
		prop := NewOrgProp()
		result, err := Run(NewModel(config, prop))
		if err != nil {
			return results, err
		}
		if result == "" {
			return results, fmt.Errorf("You must select an organization")
		}
		results.Org = result
	}

	return results, nil
}
