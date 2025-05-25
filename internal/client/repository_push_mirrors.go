package client

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"time"
)

type RepositoryPushMirror struct {
	Created       time.Time `json:"created"`
	Interval      string    `json:"interval"`
	LastError     string    `json:"last_error"`
	LastUpdate    time.Time `json:"last_update"`
	PublicKey     string    `json:"public_key"`
	RemoteAddress string    `json:"remote_address"`
	RemoteName    string    `json:"remote_name"`
	RepoName      string    `json:"repo_name"`
	SyncOnCommit  bool      `json:"sync_on_commit"`
}

func (c *Client) RepositoryPushMirrorCreate(ctx context.Context,
	owner string,
	repo string,
	interval string,
	remoteAddress string,
	remotePassword string,
	remoteUsername string,
	syncOnCommit bool,
	useSsh bool,
) (*RepositoryPushMirror, error) {
	uriRef := url.URL{Path: path.Join("api/v1/repos", owner, repo, "push_mirrors")}
	type Payload struct {
		Interval       string `json:"interval"`
		RemoteAddress  string `json:"remote_address"`
		RemotePassword string `json:"remote_password"`
		RemoteUsername string `json:"remote_username"`
		SyncOnCommit   bool   `json:"sync_on_commit"`
		UseSsh         bool   `json:"use_ssh"`
	}
	payload := Payload{
		Interval:       interval,
		RemoteAddress:  remoteAddress,
		RemotePassword: remotePassword,
		RemoteUsername: remoteUsername,
		SyncOnCommit:   syncOnCommit,
		UseSsh:         useSsh,
	}
	response := RepositoryPushMirror{}
	if _, err := c.send(ctx, "POST", &uriRef, &payload, &response); err != nil {
		return nil, fmt.Errorf("failed to create repository push mirror: %w", err)
	}
	return &response, nil
}

func (c *Client) RepositoryPushMirrorDelete(ctx context.Context, owner string, repo string, name string) (*RepositoryPushMirror, error) {
	uriRef := url.URL{Path: path.Join("api/v1/repos", owner, repo, "push_mirrors", name)}
	response := RepositoryPushMirror{}
	if _, err := c.send(ctx, "DELETE", &uriRef, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to delete repository push mirror: %w", err)
	}
	return &response, nil
}

func (c *Client) RepositoryPushMirrorGet(ctx context.Context, owner string, repo string, name string) (*RepositoryPushMirror, error) {
	uriRef := url.URL{Path: path.Join("api/v1/repos", owner, repo, "push_mirrors", name)}
	response := RepositoryPushMirror{}
	if _, err := c.send(ctx, "GET", &uriRef, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to get repository push mirrors: %w", err)
	}
	return &response, nil
}
