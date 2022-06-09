package main

import (
	"fmt"

	client "github.com/aptible/cloud-api-clients"
)

type LoginInput struct {
	email string
	pass  string
}

type CLI interface {
	Login(params *LoginInput) error
	SSH(appName string) error
	Version() string

	ListEnvironments() error
	CreateEnvironment(orgID string, params *client.EnvironmentInput) error
	RemoveEnvironment(orgID string, envID string) error

	ListApps() error
	ListDatastores() error

	ListAssetTypesForEnvironment(envID string) error
	CreateAsset(orgID string, envID string, params *client.AssetInput) error
	UpdateAsset(assetID string, envID string, ordID string, assetInput *client.AssetInput) error
	DeleteAsset(assetID string, envID string, orgID string) error
}

func main() {
	fmt.Println("Hello world!")
}
