package client

import (
	"context"
	"fmt"
	"net/url"
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

func (c *Client) OrganizationsList(ctx context.Context) ([]Organization, error) {
	var response []Organization
	uriRef := url.URL{Path: "api/v1/orgs"}
	if err := c.sendPaginated(ctx, "GET", &uriRef, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to get organizations: %w", err)
	}
	return response, nil
}
