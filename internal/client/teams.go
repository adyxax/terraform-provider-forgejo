package client

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"strconv"
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

type TeamRequest struct {
	CanCreateOrgRepo        bool              `json:"can_create_org_repo"`
	Description             string            `json:"description,omitempty"`
	IncludesAllRepositories bool              `json:"includes_all_repositories"`
	Name                    string            `json:"name"`
	Permission              string            `json:"permission"`
	Units                   []string          `json:"units"`
	UnitsMap                map[string]string `json:"units_map"`
}

func (c *Client) TeamCreate(ctx context.Context, organizationName string, payload *TeamRequest) (*Team, error) {
	uriRef := url.URL{Path: path.Join("api/v1/orgs", organizationName, "teams")}
	response := Team{}
	if _, err := c.send(ctx, "POST", &uriRef, payload, &response); err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}
	return &response, nil
}

func (c *Client) TeamDelete(ctx context.Context, id int64) error {
	uriRef := url.URL{Path: path.Join("api/v1/teams", strconv.FormatInt(id, 10))}
	response := Team{}
	if _, err := c.send(ctx, "DELETE", &uriRef, nil, &response); err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}
	return nil
}

func (c *Client) TeamGet(ctx context.Context, id int64) (*Team, error) {
	uriRef := url.URL{Path: path.Join("api/v1/teams", strconv.FormatInt(id, 10))}
	response := Team{}
	if _, err := c.send(ctx, "GET", &uriRef, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}
	return &response, nil
}

func (c *Client) TeamUpdate(ctx context.Context, id int64, payload *TeamRequest) (*Team, error) {
	uriRef := url.URL{Path: path.Join("api/v1/teams", strconv.FormatInt(id, 10))}
	response := Team{}
	if _, err := c.send(ctx, "PATCH", &uriRef, payload, &response); err != nil {
		return nil, fmt.Errorf("failed to patch team %d: %w", id, err)
	}
	return &response, nil
}

func (c *Client) TeamsList(ctx context.Context, organizationName string) ([]Team, error) {
	var response []Team
	uriRef := url.URL{Path: path.Join("api/v1/orgs", organizationName, "teams")}
	if err := c.sendPaginated(ctx, "GET", &uriRef, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to list teams of organization %s: %w", organizationName, err)
	}
	return response, nil
}
