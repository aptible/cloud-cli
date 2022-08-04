package libenv

import (
	"fmt"

	"github.com/aptible/cloud-cli/config"
	"github.com/aptible/cloud-cli/lib/org"
	"github.com/aptible/cloud-cli/ui/form"
	"github.com/charmbracelet/bubbles/list"
)

func CreateEnvOptions(orgId string) form.LoadOptionsFn {
	options := []list.Item{}
	return func(cfg *config.CloudConfig) ([]list.Item, error) {
		orgs, err := cfg.Cc.ListEnvironments(orgId)
		if err != nil {
			return options, err
		}
		for _, org := range orgs {
			options = append(options, form.FormOption{Label: org.Name, Value: org.Id})
		}
		return options, nil
	}
}

func NewEnvProp(orgId string) *form.SubSchema {
	return &form.SubSchema{
		Type:        "select",
		Title:       "Select an environment",
		LoadOptions: CreateEnvOptions(orgId),
	}
}

func EnvForm(cfg *config.CloudConfig, results *form.FormResult) error {
	if err := liborg.OrgForm(cfg, results); err != nil {
		return err
	}

	if results.Env != "" {
		return nil
	}

	prop := NewEnvProp(results.Org)
	result, err := form.Run(form.NewModel(cfg, prop))
	if err != nil {
		return err
	}
	if result == "" {
		return fmt.Errorf("You must select an environment")
	}
	results.Env = result

	return nil
}
