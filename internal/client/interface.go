package client

import (
	cloudapiclient "github.com/aptible/cloud-api-clients/clients/go"
)

/*
CloudClient
The goal of this interface is to be an abstraction layer above the cloud-api.
Whenever we want to interface with the API, we should use this interface.
*/
type CloudClient interface {
	ListEnvironments(orgId string) ([]cloudapiclient.EnvironmentOutput, int, error)
	CreateEnvironment(orgId string, params cloudapiclient.EnvironmentInput) (*cloudapiclient.EnvironmentOutput, int, error)
	DestroyEnvironment(orgId string, envId string) (int, error)

	ListOrgs() ([]cloudapiclient.OrganizationOutput, int, error)
	CreateOrg(orgId string, params cloudapiclient.OrganizationInput) (*cloudapiclient.OrganizationOutput, int, error)
	FindOrg(orgId string) (*cloudapiclient.OrganizationOutput, int, error)

	CreateAsset(orgId string, envId string, params cloudapiclient.AssetInput) (*cloudapiclient.AssetOutput, int, error)
	ListAssets(orgId string, envId string) ([]cloudapiclient.AssetOutput, int, error)
	//ListAssetTypesForEnvironment(envId string) error
	// UpdateAsset(orgId string, envId string, assetID string, assetInput *client.AssetInput) error
	DestroyAsset(orgId string, envId string, assetID string) (int, error)
}
