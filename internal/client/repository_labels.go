package client

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"strconv"
)

type RepositoryLabel struct {
	Color       string `json:"color"`
	Description string `json:"description"`
	Exclusive   bool   `json:"exclusive"`
	Id          int64  `json:"id"`
	IsArchived  bool   `json:"is_archived"`
	Name        string `json:"name"`
	Url         string `json:"url"`
}

type RepositoryLabelCreateRequest struct {
	Color       string `json:"color"`
	Description string `json:"description"`
	Exclusive   bool   `json:"exclusive"`
	IsArchived  bool   `json:"is_archived"`
	Name        string `json:"name"`
}

type RepositoryLabelUpdateRequest = RepositoryLabelCreateRequest

func (c *Client) RepositoryLabelCreate(ctx context.Context, owner string, repo string, payload *RepositoryLabelCreateRequest) (*RepositoryLabel, error) {
	uriRef := url.URL{Path: path.Join("api/v1/repos", owner, repo, "labels")}
	response := RepositoryLabel{}
	if _, err := c.send(ctx, "POST", &uriRef, payload, &response); err != nil {
		return nil, fmt.Errorf("failed to create repository label: %w", err)
	}
	return &response, nil
}

func (c *Client) RepositoryLabelDelete(ctx context.Context, owner string, repo string, id int64) error {
	uriRef := url.URL{Path: path.Join("api/v1/repos", owner, repo, "labels", strconv.Itoa(int(id)))}
	if _, err := c.send(ctx, "DELETE", &uriRef, nil, nil); err != nil {
		return fmt.Errorf("failed to delete repository label: %w", err)
	}
	return nil
}

func (c *Client) RepositoryLabelGet(ctx context.Context, owner string, repo string, id int64) (*RepositoryLabel, error) {
	uriRef := url.URL{Path: path.Join("api/v1/repos", owner, repo, "labels", strconv.Itoa(int(id)))}
	response := RepositoryLabel{}
	if _, err := c.send(ctx, "GET", &uriRef, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to get repository label: %w", err)
	}
	return &response, nil
}

func (c *Client) RepositoryLabelUpdate(ctx context.Context, owner string, repo string, id int64, payload *RepositoryLabelUpdateRequest) (*RepositoryLabel, error) {
	uriRef := url.URL{Path: path.Join("api/v1/repos", owner, repo, "labels", strconv.Itoa(int(id)))}
	response := RepositoryLabel{}
	if _, err := c.send(ctx, "PATCH", &uriRef, payload, &response); err != nil {
		return nil, fmt.Errorf("failed to update repository label: %w", err)
	}
	return &response, nil
}
