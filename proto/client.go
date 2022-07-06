package proto

import (
	client "github.com/aptible/cloud-api-clients/clients/go"
)

/*
The goal of this interface is to be an abstraction layer above the cloud-api.
Whenever we want to interface with the API, we should use this infterface.
*/
type CloudClient interface {
	ListEnvironments(orgID string) ([]client.EnvironmentOutput, error)
	CreateEnvironment(orgID string, params client.EnvironmentInput) (*client.EnvironmentOutput, error)
	DestroyEnvironment(orgID string, envID string) error

	CreateOrg(orgID string, params client.OrganizationInput) (*client.OrganizationOutput, error)
	FindOrg(orgID string) (*client.OrganizationOutput, error)

	CreateAsset(orgID string, envID string, params client.AssetInput) (*client.AssetOutput, error)

	// ListAssetTypesForEnvironment(envID string) error
	// UpdateAsset(orgID string, envID string, assetID string, assetInput *client.AssetInput) error
	// DeleteAsset(orgID string, envID string, assetID string) error
}
