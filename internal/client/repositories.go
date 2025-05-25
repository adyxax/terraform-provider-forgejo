package client

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type RepositoryExternalTracker struct {
	Description   string `json:"description"`
	Format        string `json:"external_tracker_format"`
	RegexpPattern string `json:"external_tracker_regexp_pattern"`
	Style         string `json:"external_tracker_style"`
	Url           string `json:"external_tracker_url"`
}

type RepositoryExternalWiki struct {
	Description string `json:"description"`
	Url         string `json:"external_wiki_url"`
}

type RepositoryInternalTracker struct {
	AllowOnlyContributorsToTrackTime bool `json:"allow_only_contributors_to_track_time"`
	EnableIssueDependencies          bool `json:"enable_issue_dependencies"`
	EnableTimeTracker                bool `json:"enable_time_tracker"`
}

type RepositoryTransfer struct {
	Description string             `json:"description"`
	Doer        *User              `json:"doer"`
	Recipient   *User              `json:"recipient"`
	Teams       []OrganizationTeam `json:"teams"`
}

type Repository struct {
	AllowFastForwardOnlyMerge     bool                       `json:"allow_fast_forward_only_merge"`
	AllowMergeCommits             bool                       `json:"allow_merge_commits"`
	AllowRebase                   bool                       `json:"allow_rebase"`
	AllowRebaseExplicit           bool                       `json:"allow_rebase_explicit"`
	AllowRebaseUpdate             bool                       `json:"allow_rebase_update"`
	AllowSquashMerge              bool                       `json:"allow_squash_merge"`
	ArchivedAt                    time.Time                  `json:"archived_at"`
	Archived                      bool                       `json:"archived"`
	AvatarUrl                     string                     `json:"avatar_url"`
	CloneUrl                      string                     `json:"clone_url"`
	CreatedAt                     time.Time                  `json:"created_at"`
	DefaultAllowMaintainerEdit    bool                       `json:"default_allow_maintainer_edit"`
	DefaultBranch                 string                     `json:"default_branch"`
	DefaultDeleteBranchAfterMerge bool                       `json:"default_delete_branch_after_merge"`
	DefaultMergeStyle             string                     `json:"default_merge_style"`
	DefaultUpdateStyle            string                     `json:"default_update_style"`
	Description                   string                     `json:"description"`
	Empty                         bool                       `json:"empty"`
	ExternalTracker               *RepositoryExternalTracker `json:"external_tracker"`
	ExternalWiki                  *RepositoryExternalWiki    `json:"external_wiki"`
	Fork                          bool                       `json:"fork"`
	ForksCount                    int64                      `json:"forks_count"`
	FullName                      string                     `json:"full_name"`
	GloballyEditableWiki          bool                       `json:"globally_editable_wiki"`
	HasActions                    bool                       `json:"has_actions"`
	HasIssues                     bool                       `json:"has_issues"`
	HasPackages                   bool                       `json:"has_packages"`
	HasProjects                   bool                       `json:"has_projects"`
	HasPullRequests               bool                       `json:"has_pull_requests"`
	HasReleases                   bool                       `json:"has_releases"`
	HasWiki                       bool                       `json:"has_wiki"`
	HtmlUrl                       string                     `json:"html_url"`
	Id                            int64                      `json:"id"`
	IgnoreWhitespaceConflicts     bool                       `json:"ignore_whitespace_conflicts"`
	Internal                      bool                       `json:"internal"`
	InternalTracker               *RepositoryInternalTracker `json:"internal_tracker"`
	Language                      string                     `json:"language"`
	LanguagesUrl                  string                     `json:"languages_url"`
	Link                          string                     `json:"link"`
	Mirror                        bool                       `json:"mirror"`
	MirrorInterval                string                     `json:"mirror_interval"`
	MirrorUpdated                 time.Time                  `json:"mirror_updated"`
	Name                          string                     `json:"name"`
	ObjectFormatName              string                     `json:"object_format_name"`
	OpenIssuesCount               int64                      `json:"open_issues_count"`
	OpenPrCounter                 int64                      `json:"open_pr_counter"`
	OriginalUrl                   string                     `json:"original_url"`
	Owner                         *User                      `json:"owner"`
	Parent                        *Repository                `json:"parent"`
	Permissions                   *Permission                `json:"permissions"`
	Private                       bool                       `json:"private"`
	ReleaseCounter                int64                      `json:"release_counter"`
	RepoTransfer                  *RepositoryTransfer        `json:"repo_transfer"`
	Size                          int64                      `json:"size"`
	SshUrl                        string                     `json:"ssh_url"`
	StarsCount                    int64                      `json:"stars_count"`
	Template                      bool                       `json:"template"`
	Topics                        []string                   `json:"topics"`
	UpdatedAt                     time.Time                  `json:"updated_at"`
	Url                           string                     `json:"url"`
	WatchersCount                 int64                      `json:"watchers_count"`
	Website                       string                     `json:"website"`
	WikiBranch                    string                     `json:"wiki_branch"`
}

func (c *Client) RepositoriesList(ctx context.Context) ([]Repository, error) {
	type Response struct {
		Data []Repository `json:"data"`
		Ok   bool         `json:"ok"`
	}
	uriRef := url.URL{Path: "api/v1/repos/search"}
	query := make(url.Values)
	query.Set("limit", c.maxItemsPerPageStr)
	page := 1
	var repositories []Repository
	var response Response
	for {
		query.Set("page", strconv.Itoa(page))
		uriRef.RawQuery = query.Encode()
		count, err := c.send(ctx, "GET", &uriRef, nil, &response)
		if err != nil {
			return nil, fmt.Errorf("failed to search repositories: %w", err)
		}
		if !response.Ok {
			return nil, fmt.Errorf("got a non OK status when searching repositories")
		}
		repositories = append(repositories, response.Data...)
		if count <= page*c.maxItemsPerPage {
			return repositories, nil
		}
		page++
	}
}
