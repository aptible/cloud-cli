package form

import (
	"fmt"

	"github.com/aptible/cloud-cli/internal/common"
	"github.com/charmbracelet/bubbles/list"
)

func NewEnvProp(orgId string) *SubSchema {
	return &SubSchema{
		Type:        "select",
		Title:       "Select an environment",
		LoadOptions: CreateEnvOptions(orgId),
	}
}

func EnvForm(config *common.CloudConfig, orgId string, envId string) (FormResult, error) {
	results := FormResult{
		Org: orgId,
		Env: envId,
	}

	results, err := OrgForm(config, results.Org)
	if err != nil {
		return results, err
	}

	if results.Env == "" {
		prop := NewEnvProp(results.Org)
		result, err := Run(NewModel(config, prop))
		if err != nil {
			return results, err
		}
		if result == "" {
			return results, fmt.Errorf("You must select an environment")
		}
		results.Env = result
	}

	return results, nil
}

func CreateEnvOptions(orgId string) LoadOptionsFn {
	var options []list.Item
	return func(config *common.CloudConfig) ([]list.Item, error) {
		orgs, err := config.Cc.ListEnvironments(orgId)
		if err != nil {
			return options, err
		}
		for _, org := range orgs {
			options = append(options, FormOption{Label: org.Name, Value: org.Id})
		}
		return options, nil
	}
}
