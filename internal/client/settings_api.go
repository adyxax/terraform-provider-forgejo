package client

import (
	"context"
	"fmt"
	"net/url"
)

type SettingsApi struct {
	MaxResponseItems       int `json:"max_response_items"`
	DefaultPagingNum       int `json:"default_paging_num"`
	DefaultGitTreesPerPage int `json:"default_git_trees_per_page"`
	DefaultMaxBlobSize     int `json:"default_max_blob_size"`
}

func (c *Client) settingsApiGet(ctx context.Context) (*SettingsApi, error) {
	uriRef := url.URL{Path: "api/v1/settings/api"}
	response := SettingsApi{}
	if _, err := c.send(ctx, "GET", &uriRef, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to get settings api: %w", err)
	}
	return &response, nil
}
