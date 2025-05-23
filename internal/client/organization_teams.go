package client

import (
	"context"
	"fmt"
	"net/url"
	"path"
)

type OrganizationTeam struct {
	CanCreateOrgRepo        bool              `json:"can_create_org_repo"`
	Description             string            `json:"description"`
	Id                      int64             `json:"id"`
	IncludesAllRepositories bool              `json:"includes_all_repositories"`
	Name                    string            `json:"name"`
	Organization            *Organization     `json:"organization"`
	Permission              string            `json:"permission"`
	Units                   []string          `json:"units"`
	UnitsMap                map[string]string `json:"units_map"`
}

func (c *Client) OrganizationTeamsList(ctx context.Context, organizationName string) ([]OrganizationTeam, error) {
	var response []OrganizationTeam
	uriRef := url.URL{Path: path.Join("api/v1/orgs", organizationName, "teams")}
	if err := c.sendPaginated(ctx, "GET", &uriRef, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to list teams of organization %s: %w", organizationName, err)
	}
	return response, nil
}
