package client

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"time"
)

type RepositoryActionsSecret struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
}

func (c *Client) RepositoryActionsSecretCreateOrUpdate(ctx context.Context, owner string, repo string, name string, data string) error {
	uriRef := url.URL{Path: path.Join("api/v1/repos", owner, repo, "actions/secrets", name)}
	type Payload struct {
		Data string `json:"data"`
	}
	payload := Payload{Data: data}
	if _, err := c.send(ctx, "PUT", &uriRef, &payload, nil); err != nil {
		return fmt.Errorf("failed to create or update repository actions secret: %w", err)
	}
	return nil
}

func (c *Client) RepositoryActionsSecretDelete(ctx context.Context, owner string, repo string, name string) error {
	uriRef := url.URL{Path: path.Join("api/v1/repos", owner, repo, "actions/secrets", name)}
	if _, err := c.send(ctx, "DELETE", &uriRef, nil, nil); err != nil {
		return fmt.Errorf("failed to delete repository actions secret: %w", err)
	}
	return nil
}

func (c *Client) RepositoryActionsSecretsList(ctx context.Context, owner string, repo string) ([]RepositoryActionsSecret, error) {
	uriRef := url.URL{Path: path.Join("api/v1/repos", owner, repo, "actions/secrets")}
	var response []RepositoryActionsSecret
	if err := c.sendPaginated(ctx, "GET", &uriRef, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to list repository actions secrets: %w", err)
	}
	return response, nil
}
