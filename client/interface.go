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
	DestroyEnvironment(orgId, envId string) error

	ListOrgs() ([]cac.OrganizationOutput, error)
	CreateOrg(orgId string, params cac.OrganizationInput) (*cac.OrganizationOutput, error)
	FindOrg(orgId string) (*cac.OrganizationOutput, error)

	ListAssetBundles(orgId, envId string) ([]cac.AssetBundle, error)
	CreateAsset(orgId, envId string, params cac.AssetInput) (*cac.AssetOutput, error)
	ListAssets(orgId, envId string) ([]cac.AssetOutput, error)
	DescribeAsset(orgId, envId, assetId string) (*cac.AssetOutput, error)
	//ListAssetTypesForEnvironment(envId string) error
	// UpdateAsset(orgId string, envId string, assetID string, assetInput *client.AssetInput) error
	DestroyAsset(orgId, envId, assetID string) error

	ListOperationsByAsset(orgId, assetId string) ([]cac.OperationOutput, error)

	CreateConnection(orgId, envId, assetId string, params cac.ConnectionInput) (*cac.ConnectionOutput, error)
	// DescribeConnection(orgId, envId, assetId, connectionId string) (cac.ConnectionOutput, error)
	DestroyConnection(orgId, envId, assetId, connectionId string) error
}
