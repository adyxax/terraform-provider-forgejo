package provider

import (
	"context"
	"fmt"
	"strings"

	"git.adyxax.org/adyxax/terraform-provider-forgejo/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RepositoryResource struct {
	client *client.Client
}

var _ resource.Resource = &RepositoryResource{}                // Ensure provider defined types fully satisfy framework interfaces
var _ resource.ResourceWithImportState = &RepositoryResource{} // Ensure provider defined types fully satisfy framework interfaces
func NewRepositoryResource() resource.Resource {
	return &RepositoryResource{}
}

type RepositoryResourceModel struct {
	CreatedAt       timetypes.RFC3339 `tfsdk:"created_at"`
	DefaultBranch   types.String      `tfsdk:"default_branch"`
	Description     types.String      `tfsdk:"description"`
	HasActions      types.Bool        `tfsdk:"has_actions"`
	HasIssues       types.Bool        `tfsdk:"has_issues"`
	HasPackages     types.Bool        `tfsdk:"has_packages"`
	HasProjects     types.Bool        `tfsdk:"has_projects"`
	HasPullRequests types.Bool        `tfsdk:"has_pull_requests"`
	HasReleases     types.Bool        `tfsdk:"has_releases"`
	HasWiki         types.Bool        `tfsdk:"has_wiki"`
	Name            types.String      `tfsdk:"name"`
	Owner           types.String      `tfsdk:"owner"`
	Private         types.Bool        `tfsdk:"private"`
}

func (d *RepositoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository"
}

func (d *RepositoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"created_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The creation date and time.",
			},
			"default_branch": schema.StringAttribute{
				Computed:            true,
				Default:             stringdefault.StaticString("main"),
				MarkdownDescription: "Name of the default branch. Defaults to \"main\".",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "A description string.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"has_actions": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "If true, the actions unit will be enabled. If false, the actions unit will be disabled. If unset, the server default will be left as is.",
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"has_issues": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "If true, the issues unit will be enabled. If false, the issues unit will be disabled. If unset, the server default will be left as is.",
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"has_packages": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "If true, the packages unit will be enabled. If false, the packages unit will be disabled. If unset, the server default will be left as is.",
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"has_projects": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "If true, the projects unit will be enabled. If false, the projects unit will be disabled. If unset, the server default will be left as is.",
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"has_pull_requests": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "If true, the pull requests unit will be enabled. If false, the pull requests unit will be disabled. If unset, the server default will be left as is.",
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"has_releases": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "If true, the releases unit will be enabled. If false, the releases unit will be disabled. If unset, the server default will be left as is.",
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"has_wiki": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "If true, the wiki unit will be enabled. If false, the wiki unit will be disabled. If unset, the server default will be left as is.",
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the repository.",
				Required:            true,
			},
			"owner": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the organization owning this repository. A null value here means this is a repository belonging to the user whose credentials the provider was instantiated with. Defaults to null.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"private": schema.BoolAttribute{
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "If true, the repository is private. Defaults to true.",
				Optional:            true,
			},
		},
		MarkdownDescription: "Use this resource to create and manage a git repository.",
	}
}

func (d *RepositoryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	d.client, _ = req.ProviderData.(*client.Client)
}

func (d *RepositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RepositoryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	request := client.RepositoryCreateRequest{
		DefaultBranch: data.DefaultBranch.ValueString(),
		Name:          data.Name.ValueString(),
		Private:       data.Private.ValueBool(),
	}
	if !data.Description.IsUnknown() {
		request.Description = data.Description.ValueString()
	}
	var repository *client.Repository
	var err error
	if data.Owner.IsUnknown() {
		repository, err = d.client.UserRepositoryCreate(
			ctx,
			&request)
	} else {
		repository, err = d.client.OrganizationRepositoryCreate(
			ctx,
			data.Owner.ValueString(),
			&request)
	}
	if err != nil {
		resp.Diagnostics.AddError("CreateRepository", fmt.Sprintf("failed to create repository: %s", err))
		return
	}
	updateRequest := client.RepositoryUpdateRequest{
		DefaultBranch: data.DefaultBranch.ValueString(),
		Name:          data.Name.ValueString(),
		Private:       data.Private.ValueBool(),
	}
	if !data.HasActions.IsUnknown() {
		updateRequest.HasActions = data.HasActions.ValueBoolPointer()
	}
	if !data.HasIssues.IsUnknown() {
		updateRequest.HasIssues = data.HasIssues.ValueBoolPointer()
	}
	if !data.HasPackages.IsUnknown() {
		updateRequest.HasPackages = data.HasPackages.ValueBoolPointer()
	}
	if !data.HasProjects.IsUnknown() {
		updateRequest.HasProjects = data.HasProjects.ValueBoolPointer()
	}
	if !data.HasPullRequests.IsUnknown() {
		updateRequest.HasPullRequests = data.HasPullRequests.ValueBoolPointer()
	}
	if !data.HasReleases.IsUnknown() {
		updateRequest.HasReleases = data.HasReleases.ValueBoolPointer()
	}
	if !data.HasWiki.IsUnknown() {
		updateRequest.HasWiki = data.HasWiki.ValueBoolPointer()
	}
	repository, err = d.client.RepositoryUpdate(
		ctx,
		repository.Owner.Login,
		data.Name.ValueString(),
		&updateRequest)
	if err != nil {
		resp.Diagnostics.AddError("CreateRepository", fmt.Sprintf("failed to update repository: %s", err))
		return
	}
	data.CreatedAt = timetypes.NewRFC3339TimeValue(repository.CreatedAt)
	data.Description = types.StringValue(repository.Description)
	data.HasActions = types.BoolValue(repository.HasActions)
	data.HasIssues = types.BoolValue(repository.HasIssues)
	data.HasPackages = types.BoolValue(repository.HasPackages)
	data.HasProjects = types.BoolValue(repository.HasProjects)
	data.HasPullRequests = types.BoolValue(repository.HasPullRequests)
	data.HasReleases = types.BoolValue(repository.HasReleases)
	data.HasWiki = types.BoolValue(repository.HasWiki)
	data.Owner = types.StringValue(repository.Owner.Login)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *RepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RepositoryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := d.client.RepositoryDelete(
		ctx,
		data.Owner.ValueString(),
		data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("DeleteRepository", fmt.Sprintf("failed to delete Repository: %s", err))
		return
	}
}

func (r *RepositoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "/")
	if len(idParts) == 1 && idParts[0] != "" {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("owner"), r.client.AuthenticatedUser())...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[0])...)
	} else if len(idParts) == 2 && idParts[0] != "" && idParts[1] != "" {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("owner"), idParts[0])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
	} else {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with either format <repository> or format <owner>/<repository>. Got: %q", req.ID),
		)
	}
}

func (d *RepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RepositoryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	owner := data.Owner.ValueString()
	if owner == "" {
		owner = d.client.AuthenticatedUser()
	}
	repository, err := d.client.RepositoryGet(
		ctx,
		owner,
		data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ReadRepository", fmt.Sprintf("failed to get Repository: %s", err))
		return
	}
	data.CreatedAt = timetypes.NewRFC3339TimeValue(repository.CreatedAt)
	data.DefaultBranch = types.StringValue(repository.DefaultBranch)
	data.Description = types.StringValue(repository.Description)
	data.HasActions = types.BoolValue(repository.HasActions)
	data.HasIssues = types.BoolValue(repository.HasIssues)
	data.HasPackages = types.BoolValue(repository.HasPackages)
	data.HasProjects = types.BoolValue(repository.HasProjects)
	data.HasPullRequests = types.BoolValue(repository.HasPullRequests)
	data.HasReleases = types.BoolValue(repository.HasReleases)
	data.HasWiki = types.BoolValue(repository.HasWiki)
	data.Owner = types.StringValue(repository.Owner.Login)
	data.Private = types.BoolValue(repository.Private)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *RepositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plannedData RepositoryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plannedData)...)
	var stateData RepositoryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}
	request := client.RepositoryUpdateRequest{
		DefaultBranch: plannedData.DefaultBranch.ValueString(),
		Name:          plannedData.Name.ValueString(),
		Private:       plannedData.Private.ValueBool(),
	}
	if !plannedData.Description.IsUnknown() {
		request.Description = plannedData.Description.ValueString()
	}
	if !plannedData.HasActions.IsUnknown() {
		request.HasActions = plannedData.HasActions.ValueBoolPointer()
	}
	if !plannedData.HasIssues.IsUnknown() {
		request.HasIssues = plannedData.HasIssues.ValueBoolPointer()
	}
	if !plannedData.HasPackages.IsUnknown() {
		request.HasPackages = plannedData.HasPackages.ValueBoolPointer()
	}
	if !plannedData.HasProjects.IsUnknown() {
		request.HasProjects = plannedData.HasProjects.ValueBoolPointer()
	}
	if !plannedData.HasPullRequests.IsUnknown() {
		request.HasPullRequests = plannedData.HasPullRequests.ValueBoolPointer()
	}
	if !plannedData.HasReleases.IsUnknown() {
		request.HasReleases = plannedData.HasReleases.ValueBoolPointer()
	}
	if !plannedData.HasWiki.IsUnknown() {
		request.HasWiki = plannedData.HasWiki.ValueBoolPointer()
	}
	repository, err := d.client.RepositoryUpdate(
		ctx,
		stateData.Owner.ValueString(),
		stateData.Name.ValueString(),
		&request)
	if err != nil {
		resp.Diagnostics.AddError("UpdateRepository", fmt.Sprintf("failed to update Repository: %s", err))
		return
	}
	plannedData.CreatedAt = timetypes.NewRFC3339TimeValue(repository.CreatedAt)
	if plannedData.Description.IsUnknown() {
		plannedData.Description = types.StringValue(repository.Description)
	}
	plannedData.Owner = stateData.Owner
	resp.Diagnostics.Append(resp.State.Set(ctx, &plannedData)...)
}
