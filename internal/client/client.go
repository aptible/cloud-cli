package client

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	cloudapiclient "github.com/aptible/cloud-api-clients/clients/go"
)

// client - internal cloudapiclient struct used only for this service with some common configuration
type client struct {
	ctx context.Context

	apiClient *cloudapiclient.APIClient
	debug     bool
	token     string
}

// NewClient - generate a new cloud api cloud_api_client
func NewClient(debug bool, host string, token string) CloudClient {
	config := cloudapiclient.NewConfiguration()
	config.Host = host
	config.Scheme = "https"

	apiClient := cloudapiclient.NewAPIClient(config)

	ctx := context.Background()
	ctx = context.WithValue(ctx, cloudapiclient.ContextAccessToken, token)

	return &client{
		ctx:       ctx,
		apiClient: apiClient,

		debug: debug,
		token: token,
	}
}

func (c *client) PrintResponse(r *http.Response) {
	if !c.debug {
		return
	}

	log.Println("--- DEBUG ---")
	reqDump, err := httputil.DumpRequestOut(r.Request, false)
	if err != nil {
		fmt.Println(err)
	}

	log.Printf("REQUEST:\n%s", string(reqDump))

	respDump, err := httputil.DumpResponse(r, true)
	if err != nil {
		log.Println(err)
	}

	log.Printf("RESPONSE:\n%s\n", string(respDump))
}

func (c *client) ListEnvironments(orgId string) ([]cloudapiclient.EnvironmentOutput, int, error) {
	request := c.apiClient.EnvironmentsApi.GetEnvironmentsApiV1OrganizationsOrganizationIdEnvironmentsGet(c.ctx, orgId)
	env, r, err := request.Execute()
	c.PrintResponse(r)
	return env, r.StatusCode, err
}

func (c *client) CreateEnvironment(orgId string, params cloudapiclient.EnvironmentInput) (*cloudapiclient.EnvironmentOutput, int, error) {
	request := c.apiClient.EnvironmentsApi.CreateEnvironmentApiV1OrganizationsOrganizationIdEnvironmentsPost(c.ctx, orgId).EnvironmentInput(params)
	env, r, err := request.Execute()
	c.PrintResponse(r)
	return env, r.StatusCode, err
}

func (c *client) DestroyEnvironment(orgId string, envId string) (int, error) {
	_, r, err := c.apiClient.EnvironmentsApi.DeleteEnvironmentByIdApiV1OrganizationsOrganizationIdEnvironmentsEnvironmentIdDelete(c.ctx, envId, orgId).Execute()
	c.PrintResponse(r)
	return r.StatusCode, err
}

func (c *client) CreateOrg(orgId string, params cloudapiclient.OrganizationInput) (*cloudapiclient.OrganizationOutput, int, error) {
	request := c.apiClient.OrganizationsApi.PutOrganizationApiV1OrganizationsOrganizationIdPut(c.ctx, orgId).OrganizationInput(params)
	org, r, err := request.Execute()
	c.PrintResponse(r)
	return org, r.StatusCode, err
}

func (c *client) FindOrg(orgId string) (*cloudapiclient.OrganizationOutput, int, error) {
	org, r, err := c.apiClient.OrganizationsApi.GetOrganizationByIdApiV1OrganizationsOrganizationIdGet(c.ctx, orgId).Execute()
	c.PrintResponse(r)
	return org, r.StatusCode, err
}

func (c *client) CreateAsset(orgId string, envId string, params cloudapiclient.AssetInput) (*cloudapiclient.AssetOutput, int, error) {
	request := c.apiClient.AssetsApi.CreateAssetApiV1OrganizationsOrganizationIdEnvironmentsEnvironmentIdAssetsPost(c.ctx, envId, orgId).AssetInput(params)
	asset, r, err := request.Execute()
	c.PrintResponse(r)
	return asset, r.StatusCode, err
}

func (c *client) DestroyAsset(orgId string, envId string, assetId string) (int, error) {
	request := c.apiClient.AssetsApi.DeleteAssetByIdApiV1OrganizationsOrganizationIdEnvironmentsEnvironmentIdAssetsAssetIdDelete(c.ctx, assetId, envId, orgId)
	_, r, err := request.Execute()
	c.PrintResponse(r)
	return r.StatusCode, err
}

func (c *client) ListAssets(orgId string, envId string) ([]cloudapiclient.AssetOutput, int, error) {
	request := c.apiClient.AssetsApi.GetAssetsApiV1OrganizationsOrganizationIdEnvironmentsEnvironmentIdAssetsGet(c.ctx, envId, orgId)
	assets, r, err := request.Execute()
	c.PrintResponse(r)
	return assets, r.StatusCode, err
}

func (c *client) ListOrgs() ([]cloudapiclient.OrganizationOutput, int, error) {
	request := c.apiClient.OrganizationsApi.GetOrganizationsApiV1OrganizationsGet(c.ctx)
	orgs, r, err := request.Execute()
	c.PrintResponse(r)
	return orgs, r.StatusCode, err
}
