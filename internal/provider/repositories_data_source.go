package provider

import (
	"context"
	"fmt"

	"git.adyxax.org/adyxax/terraform-provider-forgejo/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RepositoriesDataSource struct {
	client *client.Client
}

var _ datasource.DataSource = &RepositoriesDataSource{} // Ensure provider defined types fully satisfy framework interfaces
func NewRepositoriesDataSource() datasource.DataSource {
	return &RepositoriesDataSource{}
}

type RepositoriesDataSourceModel struct {
	Elements []RepositoryDataSourceModel `tfsdk:"elements"`
}
type RepositoryDataSourceModel struct {
	AllowFastForwardOnlyMerge     types.Bool                                `tfsdk:"allow_fast_forward_only_merge"`
	AllowMergeCommits             types.Bool                                `tfsdk:"allow_merge_commits"`
	AllowRebase                   types.Bool                                `tfsdk:"allow_rebase"`
	AllowRebaseExplicit           types.Bool                                `tfsdk:"allow_rebase_explicit"`
	AllowRebaseUpdate             types.Bool                                `tfsdk:"allow_rebase_update"`
	AllowSquashMerge              types.Bool                                `tfsdk:"allow_squash_merge"`
	ArchivedAt                    timetypes.RFC3339                         `tfsdk:"archived_at"`
	Archived                      types.Bool                                `tfsdk:"archived"`
	AvatarUrl                     types.String                              `tfsdk:"avatar_url"`
	CloneUrl                      types.String                              `tfsdk:"clone_url"`
	CreatedAt                     timetypes.RFC3339                         `tfsdk:"created_at"`
	DefaultAllowMaintainerEdit    types.Bool                                `tfsdk:"default_allow_maintainer_edit"`
	DefaultBranch                 types.String                              `tfsdk:"default_branch"`
	DefaultDeleteBranchAfterMerge types.Bool                                `tfsdk:"default_delete_branch_after_merge"`
	DefaultMergeStyle             types.String                              `tfsdk:"default_merge_style"`
	DefaultUpdateStyle            types.String                              `tfsdk:"default_update_style"`
	Description                   types.String                              `tfsdk:"description"`
	Empty                         types.Bool                                `tfsdk:"empty"`
	ExternalTracker               *RepositoryExternalTrackerDataSourceModel `tfsdk:"external_tracker"`
	ExternalWiki                  *RepositoryExternalWikiDataSourceModel    `tfsdk:"external_wiki"`
	Fork                          types.Bool                                `tfsdk:"fork"`
	ForksCount                    types.Int64                               `tfsdk:"forks_count"`
	FullName                      types.String                              `tfsdk:"full_name"`
	GloballyEditableWiki          types.Bool                                `tfsdk:"globally_editable_wiki"`
	HasActions                    types.Bool                                `tfsdk:"has_actions"`
	HasIssues                     types.Bool                                `tfsdk:"has_issues"`
	HasPackages                   types.Bool                                `tfsdk:"has_packages"`
	HasProjects                   types.Bool                                `tfsdk:"has_projects"`
	HasPullRequests               types.Bool                                `tfsdk:"has_pull_requests"`
	HasReleases                   types.Bool                                `tfsdk:"has_releases"`
	HasWiki                       types.Bool                                `tfsdk:"has_wiki"`
	HtmlUrl                       types.String                              `tfsdk:"html_url"`
	Id                            types.Int64                               `tfsdk:"id"`
	IgnoreWhitespaceConflicts     types.Bool                                `tfsdk:"ignore_whitespace_conflicts"`
	Internal                      types.Bool                                `tfsdk:"internal"`
	InternalTracker               *RepositoryInternalTrackerDataSourceModel `tfsdk:"internal_tracker"`
	Language                      types.String                              `tfsdk:"language"`
	LanguagesUrl                  types.String                              `tfsdk:"languages_url"`
	Link                          types.String                              `tfsdk:"link"`
	Mirror                        types.Bool                                `tfsdk:"mirror"`
	MirrorInterval                types.String                              `tfsdk:"mirror_interval"`
	MirrorUpdated                 timetypes.RFC3339                         `tfsdk:"mirror_updated"`
	Name                          types.String                              `tfsdk:"name"`
	ObjectFormatName              types.String                              `tfsdk:"object_format_name"`
	OpenIssuesCount               types.Int64                               `tfsdk:"open_issues_count"`
	OpenPrCounter                 types.Int64                               `tfsdk:"open_pr_counter"`
	OriginalUrl                   types.String                              `tfsdk:"original_url"`
	Owner                         *UserDataSourceModel                      `tfsdk:"owner"`
	//terraform does not support recursive schema definitions
	//Parent                        *RepositoryDataSourceModel                `tfsdk:"parent"`
	Permissions    *PermissionDataSourceModel         `tfsdk:"permissions"`
	Private        types.Bool                         `tfsdk:"private"`
	ReleaseCounter types.Int64                        `tfsdk:"release_counter"`
	RepoTransfer   *RepositoryTransferDataSourceModel `tfsdk:"repo_transfer"`
	Size           types.Int64                        `tfsdk:"size"`
	SshUrl         types.String                       `tfsdk:"ssh_url"`
	StarsCount     types.Int64                        `tfsdk:"stars_count"`
	Template       types.Bool                         `tfsdk:"template"`
	Topics         []types.String                     `tfsdk:"topics"`
	UpdatedAt      timetypes.RFC3339                  `tfsdk:"updated_at"`
	Url            types.String                       `tfsdk:"url"`
	WatchersCount  types.Int64                        `tfsdk:"watchers_count"`
	Website        types.String                       `tfsdk:"website"`
	WikiBranch     types.String                       `tfsdk:"wiki_branch"`
}

type RepositoryExternalTrackerDataSourceModel struct {
	Description   types.String `tfsdk:"description"`
	Format        types.String `tfsdk:"external_tracker_format"`
	RegexpPattern types.String `tfsdk:"external_tracker_regexp_pattern"`
	Style         types.String `tfsdk:"external_tracker_style"`
	Url           types.String `tfsdk:"external_tracker_url"`
}

type RepositoryExternalWikiDataSourceModel struct {
	Description types.String `tfsdk:"description"`
	Url         types.String `tfsdk:"external_wiki_url"`
}

type RepositoryInternalTrackerDataSourceModel struct {
	AllowOnlyContributorsToTrackTime types.Bool `tfsdk:"allow_only_contributors_to_track_time"`
	EnableIssueDependencies          types.Bool `tfsdk:"enable_issue_dependencies"`
	EnableTimeTracker                types.Bool `tfsdk:"enable_time_tracker"`
}

type RepositoryTransferDataSourceModel struct {
	Description types.String          `tfsdk:"description"`
	Doer        *UserDataSourceModel  `tfsdk:"doer"`
	Recipient   *UserDataSourceModel  `tfsdk:"recipient"`
	Teams       []TeamDataSourceModel `tfsdk:"teams"`
}

func (d *RepositoriesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repositories"
}

func (d *RepositoriesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"elements": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of repositories.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"allow_fast_forward_only_merge": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether fast forward only merges are allowed or not.",
						},
						"allow_merge_commits": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether merge commits are allowed or not.",
						},
						"allow_rebase": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether updating a pull request branch by rebase is allowed or not.",
						},
						"allow_rebase_explicit": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether rebase then merge commits are allowed or not.",
						},
						"allow_rebase_update": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether rebase then fast forward merges are allowed or not.",
						},
						"allow_squash_merge": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether squash merge commits are allowed on this repository or not.",
						},
						"archived_at": schema.StringAttribute{
							Computed:            true,
							CustomType:          timetypes.RFC3339Type{},
							MarkdownDescription: "The datetime at which the repository was archived.",
						},
						"archived": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the repository is archived or not.",
						},
						"avatar_url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The URL of the avatar for the repository.",
						},
						"clone_url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The URL to clone the repository.",
						},
						"created_at": schema.StringAttribute{
							Computed:            true,
							CustomType:          timetypes.RFC3339Type{},
							MarkdownDescription: "The datetime at which the repository was created.",
						},
						"default_allow_maintainer_edit": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether maintainers have edit permissions by default or not.",
						},
						"default_branch": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the default branch.",
						},
						"default_delete_branch_after_merge": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether pull request branches are deleted by default after a merge or not.",
						},
						"default_merge_style": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Name of the default merge style.",
						},
						"default_update_style": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Name of the default update style.",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "A description string.",
						},
						"empty": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the repository is empty or not.",
						},
						"external_tracker": schema.SingleNestedAttribute{
							Attributes: map[string]schema.Attribute{
								"description": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "A description string.",
								},
								"external_tracker_format": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "External issue tracker URL Format. Use the placeholders {user}, {repo} and {index} for the username, repository name and issue index.",
								},
								"external_tracker_regexp_pattern": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Regular Expression Pattern. The first captured group will be used in place of {index}.",
								},
								"external_tracker_style": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "External issue tracker Number Format.",
								},
								"external_tracker_url": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "A URL.",
								},
							},
							Computed: true,
						},
						"external_wiki": schema.SingleNestedAttribute{
							Attributes: map[string]schema.Attribute{
								"description": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "A description string.",
								},
								"external_wiki_url": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "A URL.",
								},
							},
							Computed: true,
						},
						"fork": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the repository is a fork or not.",
						},
						"forks_count": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The number of times the repository has been forked.",
						},
						"full_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The full name of the repository.",
						},
						"globally_editable_wiki": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether anyone can edit the wiki or not.",
						},
						"has_actions": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the actions unit is enabled or not.",
						},
						"has_issues": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the issues unit is enabled or not.",
						},
						"has_packages": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the packages unit is enabled or not.",
						},
						"has_projects": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the projects unit is enabled or not.",
						},
						"has_pull_requests": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the pull requests unit is enabled or not.",
						},
						"has_releases": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the releases unit is enabled or not.",
						},
						"has_wiki": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the wiki unit is enabled or not.",
						},
						"html_url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The HTTP URL of the repository.",
						},
						"id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The identifier of the repository.",
						},
						"ignore_whitespace_conflicts": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether whitespaces are ignored when detecting pull request conflicts or not.",
						},
						"internal": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether this is an internal repository or not.",
						},
						"internal_tracker": schema.SingleNestedAttribute{
							Attributes: map[string]schema.Attribute{
								"allow_only_contributors_to_track_time": schema.BoolAttribute{
									Computed:            true,
									MarkdownDescription: "Whether only contributors are allowed to track time on issues or not.",
								},
								"enable_issue_dependencies": schema.BoolAttribute{
									Computed:            true,
									MarkdownDescription: "Whether issue dependencies are enabled or not.",
								},
								"enable_time_tracker": schema.BoolAttribute{
									Computed:            true,
									MarkdownDescription: "Whether time tracking is enabled or not.",
								},
							},
							Computed: true,
						},
						"language": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The main programming language used in the repository.",
						},
						"languages_url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The URL to the languages page.",
						},
						"link": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The link.",
						},
						"mirror": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the repository is a mirror or not.",
						},
						"mirror_interval": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The mirror time interval.",
						},
						"mirror_updated": schema.StringAttribute{
							Computed:            true,
							CustomType:          timetypes.RFC3339Type{},
							MarkdownDescription: "The datetime at which the mirror was last updated.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the repository.",
						},
						"object_format_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the object format.",
						},
						"open_issues_count": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The number of open issues.",
						},
						"open_pr_counter": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The number of open pull requests.",
						},
						"original_url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The original URL.",
						},
						"owner": schema.SingleNestedAttribute{
							Attributes: userSchemaAttributes,
							Computed:   true,
						},
						//"parent"
						"permissions": schema.SingleNestedAttribute{
							Attributes: permissionSchemaAttributes,
							Computed:   true,
						},
						"private": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the repository is private or not.",
						},
						"release_counter": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The number of releases.",
						},
						"repo_transfer": schema.SingleNestedAttribute{
							Attributes: map[string]schema.Attribute{
								"description": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "A description string.",
								},
								"doer": schema.SingleNestedAttribute{
									Attributes: userSchemaAttributes,
									Computed:   true,
								},
								"recipient": schema.SingleNestedAttribute{
									Attributes: userSchemaAttributes,
									Computed:   true,
								},
								"teams": teamSchemaAttributes,
							},
							Computed: true,
						},
						"size": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The size of the repository in KiB.",
						},
						"ssh_url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The SSH URL.",
						},
						"stars_count": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The number of stars.",
						},
						"template": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the repository is a template or not.",
						},
						"topics": schema.ListAttribute{
							Computed:            true,
							ElementType:         types.StringType,
							MarkdownDescription: "The list of topics.",
						},
						"updated_at": schema.StringAttribute{
							Computed:            true,
							CustomType:          timetypes.RFC3339Type{},
							MarkdownDescription: "The datetime at which the repository was last updated.",
						},
						"url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The API URL.",
						},
						"watchers_count": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The number of watchers.",
						},
						"website": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The website URL.",
						},
						"wiki_branch": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the default branch of the wiki.",
						},
					},
				},
			},
		},
		MarkdownDescription: "Use this data source to retrieve information about existing forgejo repositories.",
	}
}

func (d *RepositoriesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client, _ = req.ProviderData.(*client.Client)
}

func (d *RepositoriesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RepositoriesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	repositories, err := d.client.RepositoriesList(ctx)
	if err != nil {
		resp.Diagnostics.AddError("ListRepositories", fmt.Sprintf("failed to list repositories: %s", err))
		return
	}
	repositoriesList := make([]RepositoryDataSourceModel, len(repositories))
	for i, repository := range repositories {
		repositoriesList[i] = RepositoryDataSourceModel{
			AllowFastForwardOnlyMerge:     types.BoolValue(repository.AllowFastForwardOnlyMerge),
			AllowMergeCommits:             types.BoolValue(repository.AllowMergeCommits),
			AllowRebase:                   types.BoolValue(repository.AllowRebase),
			AllowRebaseExplicit:           types.BoolValue(repository.AllowRebaseExplicit),
			AllowRebaseUpdate:             types.BoolValue(repository.AllowRebaseUpdate),
			AllowSquashMerge:              types.BoolValue(repository.AllowSquashMerge),
			ArchivedAt:                    timetypes.NewRFC3339TimeValue(repository.ArchivedAt),
			Archived:                      types.BoolValue(repository.Archived),
			AvatarUrl:                     types.StringValue(repository.AvatarUrl),
			CloneUrl:                      types.StringValue(repository.CloneUrl),
			CreatedAt:                     timetypes.NewRFC3339TimeValue(repository.CreatedAt),
			DefaultAllowMaintainerEdit:    types.BoolValue(repository.DefaultAllowMaintainerEdit),
			DefaultBranch:                 types.StringValue(repository.DefaultBranch),
			DefaultDeleteBranchAfterMerge: types.BoolValue(repository.DefaultDeleteBranchAfterMerge),
			DefaultMergeStyle:             types.StringValue(repository.DefaultMergeStyle),
			DefaultUpdateStyle:            types.StringValue(repository.DefaultUpdateStyle),
			Description:                   types.StringValue(repository.Description),
			Empty:                         types.BoolValue(repository.Empty),
			ExternalTracker:               nil,
			ExternalWiki:                  nil,
			Fork:                          types.BoolValue(repository.Fork),
			ForksCount:                    types.Int64Value(repository.ForksCount),
			FullName:                      types.StringValue(repository.FullName),
			GloballyEditableWiki:          types.BoolValue(repository.GloballyEditableWiki),
			HasActions:                    types.BoolValue(repository.HasActions),
			HasIssues:                     types.BoolValue(repository.HasIssues),
			HasPackages:                   types.BoolValue(repository.HasPackages),
			HasProjects:                   types.BoolValue(repository.HasProjects),
			HasPullRequests:               types.BoolValue(repository.HasPullRequests),
			HasReleases:                   types.BoolValue(repository.HasReleases),
			HasWiki:                       types.BoolValue(repository.HasWiki),
			HtmlUrl:                       types.StringValue(repository.HtmlUrl),
			Id:                            types.Int64Value(repository.Id),
			IgnoreWhitespaceConflicts:     types.BoolValue(repository.IgnoreWhitespaceConflicts),
			Internal:                      types.BoolValue(repository.Internal),
			InternalTracker:               nil,
			Language:                      types.StringValue(repository.Language),
			LanguagesUrl:                  types.StringValue(repository.LanguagesUrl),
			Link:                          types.StringValue(repository.Link),
			Mirror:                        types.BoolValue(repository.Mirror),
			MirrorInterval:                types.StringValue(repository.MirrorInterval),
			MirrorUpdated:                 timetypes.NewRFC3339TimeValue(repository.MirrorUpdated),
			Name:                          types.StringValue(repository.Name),
			ObjectFormatName:              types.StringValue(repository.ObjectFormatName),
			OpenIssuesCount:               types.Int64Value(repository.OpenIssuesCount),
			OpenPrCounter:                 types.Int64Value(repository.OpenPrCounter),
			OriginalUrl:                   types.StringValue(repository.OriginalUrl),
			Owner:                         nil,
			Permissions:                   nil,
			Private:                       types.BoolValue(repository.Private),
			ReleaseCounter:                types.Int64Value(repository.ReleaseCounter),
			RepoTransfer:                  nil,
			Size:                          types.Int64Value(repository.Size),
			SshUrl:                        types.StringValue(repository.SshUrl),
			StarsCount:                    types.Int64Value(repository.StarsCount),
			Template:                      types.BoolValue(repository.Template),
			Topics:                        make([]types.String, len(repository.Topics)),
			UpdatedAt:                     timetypes.NewRFC3339TimeValue(repository.UpdatedAt),
			Url:                           types.StringValue(repository.Url),
			WatchersCount:                 types.Int64Value(repository.WatchersCount),
			Website:                       types.StringValue(repository.Website),
			WikiBranch:                    types.StringValue(repository.WikiBranch),
		}
		if repository.ExternalTracker != nil {
			repositoriesList[i].ExternalTracker = &RepositoryExternalTrackerDataSourceModel{
				Description:   types.StringValue(repository.ExternalTracker.Description),
				Format:        types.StringValue(repository.ExternalTracker.Format),
				RegexpPattern: types.StringValue(repository.ExternalTracker.RegexpPattern),
				Style:         types.StringValue(repository.ExternalTracker.Style),
				Url:           types.StringValue(repository.ExternalTracker.Url),
			}
		}
		if repository.ExternalWiki != nil {
			repositoriesList[i].ExternalWiki = &RepositoryExternalWikiDataSourceModel{
				Description: types.StringValue(repository.ExternalTracker.Description),
				Url:         types.StringValue(repository.ExternalTracker.Url),
			}
		}
		if repository.InternalTracker != nil {
			repositoriesList[i].InternalTracker = &RepositoryInternalTrackerDataSourceModel{
				AllowOnlyContributorsToTrackTime: types.BoolValue(repository.InternalTracker.AllowOnlyContributorsToTrackTime),
				EnableIssueDependencies:          types.BoolValue(repository.InternalTracker.EnableIssueDependencies),
				EnableTimeTracker:                types.BoolValue(repository.InternalTracker.EnableTimeTracker),
			}
		}
		if repository.Owner != nil {
			repositoriesList[i].Owner = populateUserDataSourceModel(repository.Owner)
		}
		if repository.Permissions != nil {
			repositoriesList[i].Permissions = &PermissionDataSourceModel{
				Admin: types.BoolValue(repository.Permissions.Admin),
				Pull:  types.BoolValue(repository.Permissions.Pull),
				Push:  types.BoolValue(repository.Permissions.Push),
			}
		}
		if repository.RepoTransfer != nil {
			repositoriesList[i].RepoTransfer = &RepositoryTransferDataSourceModel{
				Description: types.StringValue(repository.RepoTransfer.Description),
				Doer:        populateUserDataSourceModel(repository.RepoTransfer.Doer),
				Recipient:   populateUserDataSourceModel(repository.RepoTransfer.Recipient),
				Teams:       make([]TeamDataSourceModel, len(repository.RepoTransfer.Teams)),
			}
			for j, team := range repository.RepoTransfer.Teams {
				repositoriesList[i].RepoTransfer.Teams[j] = *populateTeamDataSourceModel(&team)
			}
		}
		for j, topic := range repository.Topics {
			repositoriesList[i].Topics[j] = types.StringValue(topic)
		}
	}
	data.Elements = repositoriesList
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
