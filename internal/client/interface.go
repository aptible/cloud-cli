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
	ListEnvironments(orgId string) ([]cloudapiclient.EnvironmentOutput, error)
	CreateEnvironment(orgId string, params cloudapiclient.EnvironmentInput) (*cloudapiclient.EnvironmentOutput, error)
	DestroyEnvironment(orgId string, envId string) error

	ListOrgs() ([]cloudapiclient.OrganizationOutput, error)
	CreateOrg(orgId string, params cloudapiclient.OrganizationInput) (*cloudapiclient.OrganizationOutput, error)
	FindOrg(orgId string) (*cloudapiclient.OrganizationOutput, error)

	ListAssetBundles(orgId string, envId string) ([]cloudapiclient.AssetBundle, error)
	CreateAsset(orgId string, envId string, params cloudapiclient.AssetInput) (*cloudapiclient.AssetOutput, error)
	ListAssets(orgId string, envId string) ([]cloudapiclient.AssetOutput, error)
	DescribeAsset(orgId string, envId string, assetId string) (*cloudapiclient.AssetOutput, error)
	//ListAssetTypesForEnvironment(envId string) error
	// UpdateAsset(orgId string, envId string, assetID string, assetInput *client.AssetInput) error
	DestroyAsset(orgId string, envId string, assetID string) error

	ListOperationsByAsset(orgId string, assetId string) ([]cloudapiclient.OperationOutput, error)
}
