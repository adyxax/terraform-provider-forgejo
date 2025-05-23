package client

import (
	"context"
	"fmt"
	"net/url"
	"path"
)

type RepositoryActionsVariable struct {
	Data    string `json:"data"`
	Name    string `json:"name"`
	OwnerId int64  `json:"owner_id"`
	RepoId  int64  `json:"repo_id"`
}

func (c *Client) RepositoryActionsVariableCreate(ctx context.Context, owner string, repo string, name string, value string) error {
	uriRef := url.URL{Path: path.Join("api/v1/repos", owner, repo, "actions/variables", name)}
	type Payload struct {
		Value string `json:"value"`
	}
	payload := Payload{Value: value}
	if _, err := c.send(ctx, "POST", &uriRef, &payload, nil); err != nil {
		return fmt.Errorf("failed to create repository actions variable: %w", err)
	}
	return nil
}

func (c *Client) RepositoryActionsVariableDelete(ctx context.Context, owner string, repo string, name string) (*RepositoryActionsVariable, error) {
	uriRef := url.URL{Path: path.Join("api/v1/repos", owner, repo, "actions/variables", name)}
	response := RepositoryActionsVariable{}
	if _, err := c.send(ctx, "DELETE", &uriRef, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to delete repository actions variable: %w", err)
	}
	return &response, nil
}

func (c *Client) RepositoryActionsVariableGet(ctx context.Context, owner string, repo string, name string) (*RepositoryActionsVariable, error) {
	uriRef := url.URL{Path: path.Join("api/v1/repos", owner, repo, "actions/variables", name)}
	response := RepositoryActionsVariable{}
	if _, err := c.send(ctx, "GET", &uriRef, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to get repository actions variable: %w", err)
	}
	return &response, nil
}

func (c *Client) RepositoryActionsVariableUpdate(ctx context.Context, owner string, repo string, oldName string, newName string, value string) error {
	uriRef := url.URL{Path: path.Join("api/v1/repos", owner, repo, "actions/variables", oldName)}
	type Payload struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	payload := Payload{Name: newName, Value: value}
	if _, err := c.send(ctx, "PUT", &uriRef, &payload, nil); err != nil {
		return fmt.Errorf("failed to update repository actions variable: %w", err)
	}
	return nil
}
