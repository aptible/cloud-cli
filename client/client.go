package client

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/aptible/cloud-cli/proto"

	client "github.com/aptible/cloud-api-clients/clients/go"
)

type Client struct {
	ApiClient client.APIClient
	Ctx       context.Context
	Debug     bool
}

func NewClient(ctx context.Context, apiClient *client.APIClient, debug bool) proto.CloudClient {
	return &Client{
		ApiClient: *apiClient,
		Ctx:       ctx,
		Debug:     debug,
	}
}

func (c *Client) PrintResponse(r *http.Response) error {
	if !c.Debug {
		return nil
	}

	if r == nil {
		return fmt.Errorf("response is nil")
	}

	fmt.Println("--- DEBUG ---")
	reqDump, err := httputil.DumpRequestOut(r.Request, false)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("REQUEST:\n%s", string(reqDump))

	respDump, err := httputil.DumpResponse(r, true)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("RESPONSE:\n%s\n", string(respDump))

	return nil
}

func (c *Client) ListEnvironments(orgID string) ([]client.EnvironmentOutput, error) {
	request := c.ApiClient.EnvironmentsApi.GetEnvironmentsApiV1OrganizationsOrganizationIdEnvironmentsGet(c.Ctx, orgID)
	env, r, err := request.Execute()
	c.PrintResponse(r)
	return env, err
}

func (c *Client) CreateEnvironment(orgID string, params client.EnvironmentInput) (*client.EnvironmentOutput, error) {
	request := c.ApiClient.EnvironmentsApi.CreateEnvironmentApiV1OrganizationsOrganizationIdEnvironmentsPost(c.Ctx, orgID).EnvironmentInput(params)
	env, r, err := request.Execute()
	c.PrintResponse(r)
	return env, err
}

func (c *Client) DestroyEnvironment(orgID string, envID string) error {
	_, r, err := c.ApiClient.EnvironmentsApi.DeleteEnvironmentByIdApiV1OrganizationsOrganizationIdEnvironmentsEnvironmentIdDelete(c.Ctx, envID, orgID).Execute()
	c.PrintResponse(r)
	return err
}

func (c *Client) CreateOrg(orgID string, params client.OrganizationInput) (*client.OrganizationOutput, error) {
	request := c.ApiClient.OrganizationsApi.PutOrganizationApiV1OrganizationsOrganizationIdPut(c.Ctx, orgID).OrganizationInput(params)
	org, r, err := request.Execute()
	c.PrintResponse(r)
	return org, err
}

func (c *Client) FindOrg(orgID string) (*client.OrganizationOutput, error) {
	org, r, err := c.ApiClient.OrganizationsApi.GetOrganizationByIdApiV1OrganizationsOrganizationIdGet(c.Ctx, orgID).Execute()
	c.PrintResponse(r)
	return org, err
}

func (c *Client) ListOrgs() ([]client.OrganizationOutput, error) {
	request := c.ApiClient.OrganizationsApi.GetOrganizationsApiV1OrganizationsGet(c.Ctx)
	orgs, r, err := request.Execute()
	c.PrintResponse(r)
	return orgs, err
}

func (c *Client) CreateAsset(orgID string, envID string, params client.AssetInput) (*client.AssetOutput, error) {
	request := c.ApiClient.AssetsApi.CreateAssetApiV1OrganizationsOrganizationIdEnvironmentsEnvironmentIdAssetsPost(c.Ctx, envID, orgID).AssetInput(params)
	asset, r, err := request.Execute()
	c.PrintResponse(r)
	return asset, err
}

func (c *Client) ListAssets(orgID string, envID string) ([]client.AssetOutput, error) {
	request := c.ApiClient.AssetsApi.GetAssetsApiV1OrganizationsOrganizationIdEnvironmentsEnvironmentIdAssetsGet(c.Ctx, envID, orgID)
	assets, r, err := request.Execute()
	c.PrintResponse(r)
	return assets, err
}

func (c *Client) DestroyAsset(orgID string, envID string, assetID string) error {
	request := c.ApiClient.AssetsApi.DeleteAssetByIdApiV1OrganizationsOrganizationIdEnvironmentsEnvironmentIdAssetsAssetIdDelete(c.Ctx, envID, orgID, assetID)
	_, r, err := request.Execute()
	c.PrintResponse(r)
	return err
}
