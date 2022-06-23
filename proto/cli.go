package proto

import (
	client "github.com/aptible/cloud-api-clients/clients/go"
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
