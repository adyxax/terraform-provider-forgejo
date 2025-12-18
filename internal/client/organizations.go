package client

import (
	"context"
	"fmt"
	"net/url"
	"path"
)

type Organization struct {
	AvatarUrl                 string `json:"avatar_url"`
	Description               string `json:"description"`
	Email                     string `json:"email"`
	FullName                  string `json:"full_name"`
	Id                        int64  `json:"id"`
	Location                  string `json:"location"`
	Name                      string `json:"name"`
	RepoAdminChangeTeamAccess bool   `json:"repo_admin_change_team_access"`
	Visibility                string `json:"visibility"`
	Website                   string `json:"website"`
}

type OrganizationCreateRequest struct {
	Description               string `json:"description,omitempty"`
	Email                     string `json:"email,omitempty"`
	FullName                  string `json:"full_name,omitempty"`
	Location                  string `json:"location,omitempty"`
	RepoAdminChangeTeamAccess bool   `json:"repo_admin_change_team_access"`
	Username                  string `json:"username"`
	Visibility                string `json:"visibility,omitempty"`
	Website                   string `json:"website,omitempty"`
}

type OrganizationPatchRequest struct {
	Description               string `json:"description,omitempty"`
	Email                     string `json:"email,omitempty"`
	FullName                  string `json:"full_name,omitempty"`
	Location                  string `json:"location,omitempty"`
	RepoAdminChangeTeamAccess bool   `json:"repo_admin_change_team_access"`
	Visibility                string `json:"visibility,omitempty"`
	Website                   string `json:"website,omitempty"`
}

func (c *Client) OrganizationCreate(ctx context.Context, payload *OrganizationCreateRequest) (*Organization, error) {
	uriRef := url.URL{Path: "api/v1/orgs"}
	response := Organization{}
	if _, err := c.send(ctx, "POST", &uriRef, payload, &response); err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}
	return &response, nil
}

func (c *Client) OrganizationDelete(ctx context.Context, name string) error {
	uriRef := url.URL{Path: path.Join("api/v1/orgs", name)}
	if _, err := c.send(ctx, "DELETE", &uriRef, nil, nil); err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}
	return nil
}

func (c *Client) OrganizationGet(ctx context.Context, name string) (*Organization, error) {
	var response Organization
	uriRef := url.URL{Path: path.Join("api/v1/orgs", name)}
	if _, err := c.send(ctx, "GET", &uriRef, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	return &response, nil
}

func (c *Client) OrganizationsList(ctx context.Context) ([]Organization, error) {
	var response []Organization
	uriRef := url.URL{Path: "api/v1/orgs"}
	if err := c.sendPaginated(ctx, "GET", &uriRef, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to get organizations: %w", err)
	}
	return response, nil
}

func (c *Client) OrganizationUpdate(ctx context.Context, name string, payload *OrganizationPatchRequest) (*Organization, error) {
	uriRef := url.URL{Path: path.Join("api/v1/orgs", name)}
	response := Organization{}
	if _, err := c.send(ctx, "PATCH", &uriRef, payload, &response); err != nil {
		return nil, fmt.Errorf("failed to update organization: %w", err)
	}
	return &response, nil
}
