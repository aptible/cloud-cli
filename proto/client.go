package proto

import (
	client "github.com/aptible/cloud-api-clients"
)

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
