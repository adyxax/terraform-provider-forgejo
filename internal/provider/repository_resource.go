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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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

func populateRepositoryResourceModel(data *RepositoryResourceModel, repository *client.Repository) {
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
	data.Name = types.StringValue(repository.Name)
	data.Owner = types.StringValue(repository.Owner.Login)
	data.Private = types.BoolValue(repository.Private)
}

func repositoryUpdateRequest(data *RepositoryResourceModel) *client.RepositoryUpdateRequest {
	return &client.RepositoryUpdateRequest{
		DefaultBranch:   data.DefaultBranch.ValueString(),
		Description:     stringPointerIfKnown(data.Description),
		HasActions:      boolPointerIfKnown(data.HasActions),
		HasIssues:       boolPointerIfKnown(data.HasIssues),
		HasPackages:     boolPointerIfKnown(data.HasPackages),
		HasProjects:     boolPointerIfKnown(data.HasProjects),
		HasPullRequests: boolPointerIfKnown(data.HasPullRequests),
		HasReleases:     boolPointerIfKnown(data.HasReleases),
		HasWiki:         boolPointerIfKnown(data.HasWiki),
		Name:            data.Name.ValueString(),
		Private:         data.Private.ValueBool(),
	}
}

func (d *RepositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RepositoryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	createRequest := client.RepositoryCreateRequest{
		DefaultBranch: data.DefaultBranch.ValueString(),
		Description:   stringPointerIfKnown(data.Description),
		Name:          data.Name.ValueString(),
		Private:       data.Private.ValueBool(),
	}
	var repository *client.Repository
	var err error
	if data.Owner.IsUnknown() || data.Owner.IsNull() {
		repository, err = d.client.UserRepositoryCreate(ctx, &createRequest)
	} else {
		repository, err = d.client.OrganizationRepositoryCreate(ctx, data.Owner.ValueString(), &createRequest)
	}
	if err != nil {
		resp.Diagnostics.AddError("CreateRepository", fmt.Sprintf("failed to create repository: %s", err))
		return
	}
	// The create API does not accept all attributes, so we need a second one to
	// apply everything.
	repository, err = d.client.RepositoryUpdate(
		ctx,
		repository.Owner.Login,
		repository.Name,
		repositoryUpdateRequest(&data))
	if err != nil {
		resp.Diagnostics.AddError("CreateRepository", fmt.Sprintf("failed to update repository after creation: %s", err))
		return
	}
	populateRepositoryResourceModel(&data, repository)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *RepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RepositoryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := d.client.RepositoryDelete(ctx, data.Owner.ValueString(), data.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("DeleteRepository", fmt.Sprintf("failed to delete repository: %s", err))
		return
	}
}

func (r *RepositoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "/")
	switch {
	case len(idParts) == 1 && idParts[0] != "":
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("owner"), r.client.AuthenticatedUser())...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[0])...)
	case len(idParts) == 2 && idParts[0] != "" && idParts[1] != "":
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("owner"), idParts[0])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
	default:
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
	repository, err := d.client.RepositoryGet(ctx, owner, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ReadRepository", fmt.Sprintf("failed to get repository: %s", err))
		return
	}
	populateRepositoryResourceModel(&data, repository)
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
	repository, err := d.client.RepositoryUpdate(
		ctx,
		stateData.Owner.ValueString(),
		stateData.Name.ValueString(),
		repositoryUpdateRequest(&plannedData))
	if err != nil {
		resp.Diagnostics.AddError("UpdateRepository", fmt.Sprintf("failed to update repository: %s", err))
		return
	}
	var serverData RepositoryResourceModel
	populateRepositoryResourceModel(&serverData, repository)
	plannedData.CreatedAt = knownOr(plannedData.CreatedAt, serverData.CreatedAt)
	plannedData.DefaultBranch = knownOr(plannedData.DefaultBranch, serverData.DefaultBranch)
	plannedData.Description = knownOr(plannedData.Description, serverData.Description)
	plannedData.HasActions = knownOr(plannedData.HasActions, serverData.HasActions)
	plannedData.HasIssues = knownOr(plannedData.HasIssues, serverData.HasIssues)
	plannedData.HasPackages = knownOr(plannedData.HasPackages, serverData.HasPackages)
	plannedData.HasProjects = knownOr(plannedData.HasProjects, serverData.HasProjects)
	plannedData.HasPullRequests = knownOr(plannedData.HasPullRequests, serverData.HasPullRequests)
	plannedData.HasReleases = knownOr(plannedData.HasReleases, serverData.HasReleases)
	plannedData.HasWiki = knownOr(plannedData.HasWiki, serverData.HasWiki)
	plannedData.Name = knownOr(plannedData.Name, serverData.Name)
	plannedData.Owner = knownOr(plannedData.Owner, serverData.Owner)
	plannedData.Private = knownOr(plannedData.Private, serverData.Private)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plannedData)...)
}
