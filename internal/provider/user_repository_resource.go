package provider

import (
	"context"
	"fmt"

	"git.adyxax.org/adyxax/terraform-provider-forgejo/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type UserRepositoryResource struct {
	client *client.Client
}

var _ resource.Resource = &UserRepositoryResource{} // Ensure provider defined types fully satisfy framework interfaces
func NewUserRepositoryResource() resource.Resource {
	return &UserRepositoryResource{}
}

type UserRepositoryResourceModel struct {
	AllowFastForwardOnlyMerge     types.Bool                                `tfsdk:"allow_fast_forward_only_merge"`
	AllowManualMerge              types.Bool                                `tfsdk:"allow_manual_merge"`
	AllowMergeCommits             types.Bool                                `tfsdk:"allow_merge_commits"`
	AllowRebase                   types.Bool                                `tfsdk:"allow_rebase"`
	AllowRebaseExplicit           types.Bool                                `tfsdk:"allow_rebase_explicit"`
	AllowRebaseUpdate             types.Bool                                `tfsdk:"allow_rebase_update"`
	AllowSquashMerge              types.Bool                                `tfsdk:"allow_squash_merge"`
	ArchivedAt                    timetypes.RFC3339                         `tfsdk:"archived_at"`
	Archived                      types.Bool                                `tfsdk:"archived"`
	AutodetectManualMerge         types.Bool                                `tfsdk:"autodetect_manual_merge"`
	AutoInit                      types.Bool                                `tfsdk:"auto_init"`
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
	EnablePrune                   types.Bool                                `tfsdk:"enable_prune"`
	ExternalTracker               *RepositoryExternalTrackerDataSourceModel `tfsdk:"external_tracker"`
	ExternalWiki                  *RepositoryExternalWikiDataSourceModel    `tfsdk:"external_wiki"`
	Fork                          types.Bool                                `tfsdk:"fork"`
	ForksCount                    types.Int64                               `tfsdk:"forks_count"`
	FullName                      types.String                              `tfsdk:"full_name"`
	Gitignores                    types.String                              `tfsdk:"gitignores"`
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
	IssueLabels                   types.String                              `tfsdk:"issue_labels"`
	Language                      types.String                              `tfsdk:"language"`
	LanguagesUrl                  types.String                              `tfsdk:"languages_url"`
	License                       types.String                              `tfsdk:"license"`
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
	//Parent                        *UserRepositoryDataSourceModel                `tfsdk:"parent"`
	Permissions    *PermissionDataSourceModel         `tfsdk:"permissions"`
	Private        types.Bool                         `tfsdk:"private"`
	Readme         types.String                       `tfsdk:"readme"`
	ReleaseCounter types.Int64                        `tfsdk:"release_counter"`
	RepoTransfer   *RepositoryTransferDataSourceModel `tfsdk:"repo_transfer"`
	Size           types.Int64                        `tfsdk:"size"`
	SshUrl         types.String                       `tfsdk:"ssh_url"`
	StarsCount     types.Int64                        `tfsdk:"stars_count"`
	Template       types.Bool                         `tfsdk:"template"`
	Topics         []types.String                     `tfsdk:"topics"`
	TrustModel     types.String                       `tfsdk:"trust_model"`
	UpdatedAt      timetypes.RFC3339                  `tfsdk:"updated_at"`
	Url            types.String                       `tfsdk:"url"`
	WatchersCount  types.Int64                        `tfsdk:"watchers_count"`
	Website        types.String                       `tfsdk:"website"`
	WikiBranch     types.String                       `tfsdk:"wiki_branch"`
}

func (d *UserRepositoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_repository"
}

func (d *UserRepositoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"allow_fast_forward_only_merge": schema.BoolAttribute{
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Whether fast forward only merges are allowed or not. Defaults to true.",
				Optional:            true,
			},
			"allow_manual_merge": schema.BoolAttribute{
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Whether manual merges are allowed or not. Defaults to true.",
				Optional:            true,
			},
			"allow_merge_commits": schema.BoolAttribute{
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Whether merge commits are allowed or not. Defaults to true.",
				Optional:            true,
			},
			"allow_rebase": schema.BoolAttribute{
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Whether rebases are allowed or not. Defaults to true.",
				Optional:            true,
			},
			"allow_rebase_explicit": schema.BoolAttribute{
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Whether explicit rebases are allowed or not. Defaults to true.",
				Optional:            true,
			},
		},
		MarkdownDescription: "Use this resource to create and manage a repository belonging to the user whose credentials you are instantiating the provider with.",
	}
}

func (d *UserRepositoryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	d.client, _ = req.ProviderData.(*client.Client)
}

func (d *UserRepositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data UserRepositoryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	request := client.UserRepositoryCreateRequest{
		Color:       data.Color.ValueString(),
		Description: data.Description.ValueString(),
		Exclusive:   data.Exclusive.ValueBool(),
		IsArchived:  data.IsArchived.ValueBool(),
		Name:        data.Name.ValueString(),
	}
	label, err := d.client.UserRepositoryCreate(
		ctx,
		data.Owner.ValueString(),
		data.UserRepository.ValueString(),
		&request)
	if err != nil {
		resp.Diagnostics.AddError("CreateUserRepository", fmt.Sprintf("failed to create UserRepository: %s", err))
		return
	}
	data.Id = types.Int64Value(label.Id)
	data.Url = types.StringValue(label.Url)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *UserRepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data UserRepositoryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := d.client.UserRepositoryDelete(
		ctx,
		data.Owner.ValueString(),
		data.UserRepository.ValueString(),
		data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("DeleteUserRepository", fmt.Sprintf("failed to delete UserRepository: %s", err))
		return
	}
}

func (d *UserRepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data UserRepositoryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	label, err := d.client.UserRepositoryGet(
		ctx,
		data.Owner.ValueString(),
		data.UserRepository.ValueString(),
		data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("ReadUserRepository", fmt.Sprintf("failed to get UserRepository: %s", err))
		return
	}
	data.Color = types.StringValue(label.Color)
	data.Description = types.StringValue(label.Description)
	data.Exclusive = types.BoolValue(label.Exclusive)
	data.IsArchived = types.BoolValue(label.IsArchived)
	data.Name = types.StringValue(label.Name)
	data.Url = types.StringValue(label.Url)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *UserRepositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plannedData UserRepositoryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plannedData)...)
	var stateData UserRepositoryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}
	request := client.UserRepositoryUpdateRequest{
		Color:       plannedData.Color.ValueString(),
		Description: plannedData.Description.ValueString(),
		Exclusive:   plannedData.Exclusive.ValueBool(),
		IsArchived:  plannedData.IsArchived.ValueBool(),
		Name:        plannedData.Name.ValueString(),
	}
	label, err := d.client.UserRepositoryUpdate(
		ctx,
		plannedData.Owner.ValueString(),
		plannedData.UserRepository.ValueString(),
		stateData.Id.ValueInt64(),
		&request)
	if err != nil {
		resp.Diagnostics.AddError("UpdateUserRepository", fmt.Sprintf("failed to update UserRepository: %s", err))
		return
	}
	plannedData.Id = types.Int64Value(label.Id)
	plannedData.Url = types.StringValue(label.Url)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plannedData)...)
}
