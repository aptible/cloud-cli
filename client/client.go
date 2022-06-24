package client

import (
	"context"

	"github.com/aptible/cloud-cli/proto"

	client "github.com/aptible/cloud-api-clients/clients/go"
)

type Client struct {
	ApiClient client.APIClient
	Ctx       context.Context
}

func NewClient(ctx context.Context, apiClient *client.APIClient) proto.CloudClient {
	return &Client{
		ApiClient: *apiClient,
		Ctx:       ctx,
	}
}

func (c *Client) ListEnvironments(orgID string) ([]client.EnvironmentOutput, error) {
	env, _, err := c.ApiClient.EnvironmentsApi.GetEnvironmentsApiV1OrganizationsOrganizationIdEnvironmentsGet(c.Ctx, orgID).Execute()
	return env, err
}

func (c *Client) CreateEnvironment(orgID string, params client.EnvironmentInput) (*client.EnvironmentOutput, error) {
	env, _, err := c.ApiClient.EnvironmentsApi.CreateEnvironmentApiV1OrganizationsOrganizationIdEnvironmentsPost(c.Ctx, orgID).EnvironmentInput(params).Execute()
	return env, err
}

func (c *Client) DestroyEnvironment(orgID string, envID string) error {
	_, _, err := c.ApiClient.EnvironmentsApi.DeleteEnvironmentByIdApiV1OrganizationsOrganizationIdEnvironmentsEnvironmentIdDelete(c.Ctx, envID, orgID).Execute()
	return err
}
