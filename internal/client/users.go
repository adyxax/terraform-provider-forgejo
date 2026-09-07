package client

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"time"
)

type User struct {
	Active           bool      `json:"active"`
	AvatarUrl        string    `json:"avatar_url"`
	Created          time.Time `json:"created"`
	Description      string    `json:"description"`
	Email            string    `json:"email"`
	FollowerCount    int64     `json:"followers_count"`
	FollowingCount   int64     `json:"following_count"`
	FullName         string    `json:"full_name"`
	HtmlUrl          string    `json:"html_url"`
	Id               int64     `json:"id"`
	IsAdmin          bool      `json:"is_admin"`
	Language         string    `json:"language"`
	LastLogin        time.Time `json:"last_login"`
	Location         string    `json:"location"`
	LoginName        string    `json:"login_name"`
	Login            string    `json:"login"`
	ProhibitLogin    bool      `json:"prohibit_login"`
	Pronouns         string    `json:"pronouns"`
	Restricted       bool      `json:"restricted"`
	SourceId         int64     `json:"source_id"`
	StarredRepoCount int64     `json:"starred_repos_count"`
	Visibility       string    `json:"visibility"`
	Website          string    `json:"website"`
}

type UserCreateRequest struct {
	Email              string `json:"email"`
	FullName           string `json:"full_name,omitempty"`
	LoginName          string `json:"login_name,omitempty"`
	MustChangePassword *bool  `json:"must_change_password,omitempty"`
	Password           string `json:"password,omitempty"`
	Restricted         *bool  `json:"restricted,omitempty"`
	SendNotify         bool   `json:"send_notify,omitempty"`
	SourceId           int64  `json:"source_id,omitempty"`
	Username           string `json:"username"`
	Visibility         string `json:"visibility,omitempty"`
}

type UserUpdateRequest struct {
	Active                  *bool   `json:"active,omitempty"`
	Admin                   *bool   `json:"admin,omitempty"`
	AllowCreateOrganization *bool   `json:"allow_create_organization,omitempty"`
	AllowGitHook            *bool   `json:"allow_git_hook,omitempty"`
	AllowImportLocal        *bool   `json:"allow_import_local,omitempty"`
	Description             *string `json:"description,omitempty"`
	Email                   *string `json:"email,omitempty"`
	FullName                *string `json:"full_name,omitempty"`
	Location                *string `json:"location,omitempty"`
	LoginName               string  `json:"login_name"`
	MaxRepoCreation         *int64  `json:"max_repo_creation,omitempty"`
	MustChangePassword      *bool   `json:"must_change_password,omitempty"`
	Password                string  `json:"password,omitempty"`
	ProhibitLogin           *bool   `json:"prohibit_login,omitempty"`
	Pronouns                *string `json:"pronouns,omitempty"`
	Restricted              *bool   `json:"restricted,omitempty"`
	SourceId                int64   `json:"source_id"`
	Visibility              string  `json:"visibility,omitempty"`
	Website                 *string `json:"website,omitempty"`
}

func (c *Client) authenticatedUserGet(ctx context.Context) (*User, error) {
	uriRef := url.URL{Path: "api/v1/user"}
	var response User
	if _, err := c.send(ctx, "GET", &uriRef, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to get authenticated user: %w", err)
	}
	return &response, nil
}

func (c *Client) UserCreate(ctx context.Context, payload *UserCreateRequest) (*User, error) {
	uriRef := url.URL{Path: "api/v1/admin/users"}
	response := User{}
	if _, err := c.send(ctx, "POST", &uriRef, payload, &response); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &response, nil
}

func (c *Client) UserDelete(ctx context.Context, username string) error {
	uriRef := url.URL{Path: path.Join("api/v1/admin/users", username)}
	if _, err := c.send(ctx, "DELETE", &uriRef, nil, nil); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (c *Client) UserGet(ctx context.Context, username string) (*User, error) {
	uriRef := url.URL{Path: path.Join("api/v1/users", username)}
	response := User{}
	if _, err := c.send(ctx, "GET", &uriRef, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &response, nil
}

func (c *Client) UsersList(ctx context.Context) ([]User, error) {
	type Response struct {
		Data []User `json:"data"`
		Ok   bool   `json:"ok"`
	}
	uriRef := url.URL{Path: "api/v1/users/search"}
	query := make(url.Values)
	query.Set("limit", c.maxItemsPerPageStr)
	page := 1
	var users []User
	var response Response
	for {
		query.Set("page", strconv.Itoa(page))
		uriRef.RawQuery = query.Encode()
		count, err := c.send(ctx, "GET", &uriRef, nil, &response)
		if err != nil {
			return nil, fmt.Errorf("failed to search users: %w", err)
		}
		if !response.Ok {
			return nil, fmt.Errorf("got a non OK status when searching users")
		}
		users = append(users, response.Data...)
		if count <= page*c.maxItemsPerPage {
			return users, nil
		}
		page++
	}
}

func (c *Client) UserUpdate(ctx context.Context, username string, payload *UserUpdateRequest) (*User, error) {
	uriRef := url.URL{Path: path.Join("api/v1/admin/users", username)}
	response := User{}
	if _, err := c.send(ctx, "PATCH", &uriRef, payload, &response); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return &response, nil
}
