package client

import (
	"context"
	"fmt"
	"net/url"
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
