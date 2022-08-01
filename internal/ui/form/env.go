package form

import (
	"fmt"

	"github.com/aptible/cloud-cli/internal/common"
	"github.com/charmbracelet/bubbles/list"
)

func CreateEnvOptions(orgId string) LoadOptionsFn {
	options := []list.Item{}
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

func NewEnvProp(orgId string) *SubSchema {
	return &SubSchema{
		Type:        "select",
		Title:       "Select an environment",
		LoadOptions: CreateEnvOptions(orgId),
	}
}

func EnvForm(config *common.CloudConfig, results *FormResult) error {
	if err := OrgForm(config, results); err != nil {
		return err
	}

	if results.Env != "" {
		return nil
	}

	prop := NewEnvProp(results.Org)
	result, err := Run(NewModel(config, prop))
	if err != nil {
		return err
	}
	if result == "" {
		return fmt.Errorf("You must select an environment")
	}
	results.Env = result

	return nil
}
