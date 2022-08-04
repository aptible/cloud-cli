package client

import (
	cac "github.com/aptible/cloud-api-clients/clients/go"
)

/*
CloudClient
The goal of this interface is to be an abstraction layer above the cloud-api.
Whenever we want to interface with the API, we should use this interface.
*/
type CloudClient interface {
	ListEnvironments(orgId string) ([]cac.EnvironmentOutput, error)
	CreateEnvironment(orgId string, params cac.EnvironmentInput) (*cac.EnvironmentOutput, error)
	DestroyEnvironment(orgId string, envId string) error

	ListOrgs() ([]cac.OrganizationOutput, error)
	CreateOrg(orgId string, params cac.OrganizationInput) (*cac.OrganizationOutput, error)
	FindOrg(orgId string) (*cac.OrganizationOutput, error)

	ListAssetBundles(orgId string, envId string) ([]cac.AssetBundle, error)
	CreateAsset(orgId string, envId string, params cac.AssetInput) (*cac.AssetOutput, error)
	ListAssets(orgId string, envId string) ([]cac.AssetOutput, error)
	DescribeAsset(orgId string, envId string, assetId string) (*cac.AssetOutput, error)
	//ListAssetTypesForEnvironment(envId string) error
	// UpdateAsset(orgId string, envId string, assetID string, assetInput *client.AssetInput) error
	DestroyAsset(orgId string, envId string, assetID string) error

	ListOperationsByAsset(orgId string, assetId string) ([]cac.OperationOutput, error)
}
