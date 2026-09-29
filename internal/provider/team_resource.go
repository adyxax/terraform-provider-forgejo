package provider

import (
	"context"
	"fmt"
	"strconv"

	"git.adyxax.org/adyxax/terraform-provider-forgejo/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TeamResource struct {
	client *client.Client
}

var _ resource.Resource = &TeamResource{}                // Ensure provider defined types fully satisfy framework interfaces
var _ resource.ResourceWithImportState = &TeamResource{} // Ensure provider defined types fully satisfy framework interfaces
func NewTeamResource() resource.Resource {
	return &TeamResource{}
}

type TeamResourceModel struct {
	CanCreateOrgRepo        types.Bool   `tfsdk:"can_create_org_repo"`
	Description             types.String `tfsdk:"description"`
	Id                      types.Int64  `tfsdk:"id"`
	IncludesAllRepositories types.Bool   `tfsdk:"includes_all_repositories"`
	Name                    types.String `tfsdk:"name"`
	OrganizationName        types.String `tfsdk:"organization_name"`
	Permission              types.String `tfsdk:"permission"`
}

func (d *TeamResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (d *TeamResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"can_create_org_repo": schema.BoolAttribute{
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether members of this team can create repositories that will belong to the organization. Defaults to false.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "A description string. If unset, the server value will be left as is.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The identifier of the team.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"includes_all_repositories": schema.BoolAttribute{
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether members of this team can access all the repositories that belong to the organization. Defaults to false.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the team.",
				Required:            true,
			},
			"organization_name": schema.StringAttribute{
				MarkdownDescription: "The name of the organization the team is a part of.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Required: true,
			},
			"permission": schema.StringAttribute{
				MarkdownDescription: "The members' permission level on the organization. Valid values are `admin`, `read` and `write`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("admin", "read", "write"),
				},
			},
		},
		MarkdownDescription: "Use this resource to create and manage a team.",
	}
}

func (d *TeamResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	d.client, _ = req.ProviderData.(*client.Client)
}

func teamUnits(permission string) []string {
	// admin teams implicitly get all units.
	if permission == "admin" {
		return nil
	}
	return []string{
		"repo.code",
		"repo.issues",
		"repo.ext_issues",
		"repo.releases",
		"repo.wiki",
		"repo.ext_wiki",
		"repo.packages",
		"repo.pulls",
		"repo.projects",
		"repo.actions",
	}
}

func populateTeamResourceModel(data *TeamResourceModel, team *client.Team) {
	data.CanCreateOrgRepo = types.BoolValue(team.CanCreateOrgRepo)
	data.Description = types.StringValue(team.Description)
	data.Id = types.Int64Value(team.Id)
	data.IncludesAllRepositories = types.BoolValue(team.IncludesAllRepositories)
	data.Name = types.StringValue(team.Name)
	// forgejo does not set team.Organization on PATCH API responses
	if team.Organization != nil {
		data.OrganizationName = types.StringValue(team.Organization.Name)
	}
	data.Permission = types.StringValue(team.Permission)
}

func teamRequest(data *TeamResourceModel) *client.TeamRequest {
	return &client.TeamRequest{
		CanCreateOrgRepo:        data.CanCreateOrgRepo.ValueBool(),
		Description:             stringPointerIfKnown(data.Description),
		IncludesAllRepositories: data.IncludesAllRepositories.ValueBool(),
		Name:                    data.Name.ValueString(),
		Permission:              data.Permission.ValueString(),
		Units:                   teamUnits(data.Permission.ValueString()),
	}
}

func (d *TeamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TeamResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	team, err := d.client.TeamCreate(ctx, data.OrganizationName.ValueString(), teamRequest(&data))
	if err != nil {
		resp.Diagnostics.AddError("CreateTeam", fmt.Sprintf("failed to create team: %s", err))
		return
	}
	populateTeamResourceModel(&data, team)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *TeamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TeamResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := d.client.TeamDelete(ctx, data.Id.ValueInt64()); err != nil {
		resp.Diagnostics.AddError("DeleteTeam", fmt.Sprintf("failed to delete team: %s", err))
		return
	}
}

func (r *TeamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("ImportTeam", fmt.Sprintf("Expected import identifier to be the team's numeric id: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (d *TeamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TeamResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	team, err := d.client.TeamGet(ctx, data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("ReadTeam", fmt.Sprintf("failed to get team: %s", err))
		return
	}
	populateTeamResourceModel(&data, team)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *TeamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plannedData TeamResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plannedData)...)
	var stateData TeamResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}
	team, err := d.client.TeamUpdate(ctx, stateData.Id.ValueInt64(), teamRequest(&plannedData))
	if err != nil {
		resp.Diagnostics.AddError("UpdateTeam", fmt.Sprintf("failed to update team: %s", err))
		return
	}
	var serverData TeamResourceModel
	populateTeamResourceModel(&serverData, team)
	plannedData.CanCreateOrgRepo = knownOr(plannedData.CanCreateOrgRepo, serverData.CanCreateOrgRepo)
	plannedData.Description = knownOr(plannedData.Description, serverData.Description)
	plannedData.Id = knownOr(plannedData.Id, serverData.Id)
	plannedData.IncludesAllRepositories = knownOr(plannedData.IncludesAllRepositories, serverData.IncludesAllRepositories)
	plannedData.Name = knownOr(plannedData.Name, serverData.Name)
	// organization_name is Required and RequiresReplace so it is always known
	plannedData.Permission = knownOr(plannedData.Permission, serverData.Permission)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plannedData)...)
}
