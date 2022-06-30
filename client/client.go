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

	fmt.Println(r.Request.Body)
	reqDump, err := httputil.DumpRequestOut(r.Request, true)
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
	env, r, err := c.ApiClient.EnvironmentsApi.GetEnvironmentsApiV1OrganizationsOrganizationIdEnvironmentsGet(c.Ctx, orgID).Execute()
	c.PrintResponse(r)
	return env, err
}

func (c *Client) CreateEnvironment(orgID string, params client.EnvironmentInput) (*client.EnvironmentOutput, error) {
	env, r, err := c.ApiClient.EnvironmentsApi.CreateEnvironmentApiV1OrganizationsOrganizationIdEnvironmentsPost(c.Ctx, orgID).EnvironmentInput(params).Execute()
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
