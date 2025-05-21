package client

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"time"
)

type Secret struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
}

func (c *Client) RepoActionSecretDelete(ctx context.Context, owner string, repo string, name string) error {
	uriRef := url.URL{Path: path.Join("api/v1/repos", owner, repo, "actions/secrets", name)}
	if _, err := c.send(ctx, "DELETE", &uriRef, nil, nil); err != nil {
		return fmt.Errorf("failed to delete repository secret: %w", err)
	}
	return nil
}

func (c *Client) RepoActionSecretPut(ctx context.Context, owner string, repo string, name string, data string) error {
	uriRef := url.URL{Path: path.Join("api/v1/repos", owner, repo, "actions/secrets", name)}
	type Payload struct {
		Data string `json:"data"`
	}
	payload := Payload{Data: data}
	if _, err := c.send(ctx, "PUT", &uriRef, &payload, nil); err != nil {
		return fmt.Errorf("failed to put repository secret: %w", err)
	}
	return nil
}

func (c *Client) RepoActionSecretsList(ctx context.Context, owner string, repo string) ([]Secret, error) {
	uriRef := url.URL{Path: path.Join("api/v1/repos", owner, repo, "actions/secrets")}
	var response []Secret
	if err := c.sendPaginated(ctx, "GET", &uriRef, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to list repository secrets: %w", err)
	}
	return response, nil
}
