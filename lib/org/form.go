package liborg

import (
	"fmt"

	"github.com/aptible/cloud-cli/config"
	"github.com/aptible/cloud-cli/ui/form"
	"github.com/charmbracelet/bubbles/list"
)

func NewOrgProp() *form.SubSchema {
	return &form.SubSchema{
		Type:        "select",
		Title:       "Select an organization",
		LoadOptions: CreateOrgOptions(),
	}
}

func CreateOrgOptions() form.LoadOptionsFn {
	options := []list.Item{}
	return func(cfg *config.CloudConfig) ([]list.Item, error) {
		orgs, err := cfg.Cc.ListOrgs()
		if err != nil {
			return options, err
		}
		for _, org := range orgs {
			options = append(options, form.FormOption{Label: org.Name, Value: org.Id})
		}
		return options, nil
	}
}

func OrgForm(cfg *config.CloudConfig, results *form.FormResult) error {
	if results.Org != "" {
		return nil
	}

	prop := NewOrgProp()
	result, err := form.Run(form.NewModel(cfg, prop))
	if err != nil {
		return err
	}
	if result == "" {
		return fmt.Errorf("You must select an organization")
	}
	results.Org = result

	return nil
}
