package client

import (
	"context"
	"fmt"
	"net/url"
	"path"
)

type Team struct {
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

func (c *Client) TeamsList(ctx context.Context, organizationName string) ([]Team, error) {
	var response []Team
	query := make(url.Values)
	query.Set("limit", "50")
	query.Set("page", "1")
	uriRef := url.URL{
		Path:     path.Join("api/v1/orgs", organizationName, "teams"),
		RawQuery: query.Encode(),
	}
	if err := c.Send(ctx, "GET", &uriRef, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to list teams of organization %s: %w", organizationName, err)
	}
	return response, nil
}
