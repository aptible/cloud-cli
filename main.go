package main

import (
	client "github.com/aptible/cloud-api-clients"
	"github.com/aptible/cloud-cli/cmd"
)

type LoginInput struct {
	email string
	pass  string
}

type CreateAssetParams struct {
	Name          string // e.g. name of db
	Engine        string // e.g. postgres
	EngineVersion string // e.g. 14.2
}

/*
The goal of this interface is to be an abstraction layer above the cloud-api.
Whenever we want to interface with the API, we should use this infterface.
*/
type CloudClient interface {
	ListEnvironments() error
	CreateEnvironment(orgID string, params *client.EnvironmentInput) error
	RemoveEnvironment(orgID string, envID string) error

	ListAssetTypesForEnvironment(envID string) error
	CreateAsset(orgID string, envID string, params *client.AssetInput) error
	UpdateAsset(orgID string, envID string, assetID string, assetInput *client.AssetInput) error
	DeleteAsset(orgID string, envID string, assetID string) error
}

/*
The goal of this interface is to represent the commands that we plan to implement.
Whatever CLI framework we use will primarily interact with this interface.
*/
type CLI interface {
	// Login(params *LoginInput) error
	// Version() string

	ListOrgs() error
	// ListApps(orgID string) error
	// ListsLogs(orgID string) error
	// ListLogsForAsset(orgID string, envID string, assetID string) error

	// SSH(orgID string, envID string, assetID string) error
	// Info(orgID string, envID string, assetID string) error
	// Open(orgID string, envID string, assetID string) error

	ListEnvironments(orgID string) error
	CreateEnvironment(orgID string, params *client.EnvironmentInput) error
	RemoveEnvironment(orgID string, envID string) error

	CreateDatastore(input *CreateAssetParams) error
	ListDatastores(orgID string) error

	// CreateBackup(orgID string, envID string, assetID string) error
	// DeleteDatastore(orgID string, envID string, assetID string) error
	// DeleteBackup(orgID string, envID string, backupID string) error
	// RestoreDatastore(backupID string, orgID string, envID string, assetID string) error
}

func main() {
	cmd.Execute()
}
