package provider

import (
	"context"
	"fmt"

	"git.adyxax.org/adyxax/terraform-provider-forgejo/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type OrganizationResource struct {
	client *client.Client
}

var _ resource.Resource = &OrganizationResource{}                // Ensure provider defined types fully satisfy framework interfaces
var _ resource.ResourceWithImportState = &OrganizationResource{} // Ensure provider defined types fully satisfy framework interfaces
func NewOrganizationResource() resource.Resource {
	return &OrganizationResource{}
}

type OrganizationResourceModel struct {
	Description               types.String `tfsdk:"description"`
	Email                     types.String `tfsdk:"email"`
	FullName                  types.String `tfsdk:"full_name"`
	Id                        types.Int64  `tfsdk:"id"`
	Location                  types.String `tfsdk:"location"`
	Name                      types.String `tfsdk:"name"`
	RepoAdminChangeTeamAccess types.Bool   `tfsdk:"repo_admin_change_team_access"`
	Visibility                types.String `tfsdk:"visibility"`
	Website                   types.String `tfsdk:"website"`
}

func (d *OrganizationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (d *OrganizationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "A description string.",
				Optional:            true,
			},
			"email": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The organization's email address.",
				Optional:            true,
			},
			"full_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The organization's full name.",
				Optional:            true,
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The identifier of the organization.",
			},
			"location": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The organization's advertised location.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the organization.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Required: true,
			},
			"repo_admin_change_team_access": schema.BoolAttribute{
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Whether an admin of a repository that belongs to this organization can change team access or not. Defaults to true.",
				Optional:            true,
			},
			"visibility": schema.StringAttribute{
				Computed:            true,
				Default:             stringdefault.StaticString("private"),
				MarkdownDescription: "The organization's visibility option: `limited`, `private`, `public`. Defaults to `private`.",
				Optional:            true,
			},
			"website": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The organization's advertised website.",
				Optional:            true,
			},
		},
		MarkdownDescription: "Use this resource to create and manage an organization.",
	}
}

func (d *OrganizationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	d.client, _ = req.ProviderData.(*client.Client)
}

func (d *OrganizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OrganizationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	request := client.OrganizationCreateRequest{
		RepoAdminChangeTeamAccess: data.RepoAdminChangeTeamAccess.ValueBool(),
		Username:                  data.Name.ValueString(),
	}
	if !data.Description.IsUnknown() {
		request.Description = data.Description.ValueString()
	}
	if !data.Email.IsUnknown() {
		request.Email = data.Email.ValueString()
	}
	if !data.FullName.IsUnknown() {
		request.FullName = data.FullName.ValueString()
	}
	if !data.Location.IsUnknown() {
		request.Location = data.Location.ValueString()
	}
	if !data.Visibility.IsUnknown() {
		request.Visibility = data.Visibility.ValueString()
	}
	if !data.Website.IsUnknown() {
		request.Website = data.Website.ValueString()
	}
	organization, err := d.client.OrganizationCreate(
		ctx,
		&request)
	if err != nil {
		resp.Diagnostics.AddError("CreateOrganization", fmt.Sprintf("failed to create organization: %s", err))
		return
	}
	data.Description = types.StringValue(organization.Description)
	data.Email = types.StringValue(organization.Email)
	data.FullName = types.StringValue(organization.FullName)
	data.Id = types.Int64Value(organization.Id)
	data.Location = types.StringValue(organization.Location)
	data.RepoAdminChangeTeamAccess = types.BoolValue(organization.RepoAdminChangeTeamAccess)
	data.Visibility = types.StringValue(organization.Visibility)
	data.Website = types.StringValue(organization.Website)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *OrganizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OrganizationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := d.client.OrganizationDelete(
		ctx,
		data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("DeleteOrganization", fmt.Sprintf("failed to delete organization: %s", err))
		return
	}
}

func (r *OrganizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), req.ID)...)
}

func (d *OrganizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OrganizationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	organization, err := d.client.OrganizationGet(
		ctx,
		data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ReadOrganization", fmt.Sprintf("failed to get organization: %s", err))
		return
	}
	data.Description = types.StringValue(organization.Description)
	data.Email = types.StringValue(organization.Email)
	data.FullName = types.StringValue(organization.FullName)
	data.Id = types.Int64Value(organization.Id)
	data.Location = types.StringValue(organization.Location)
	data.RepoAdminChangeTeamAccess = types.BoolValue(organization.RepoAdminChangeTeamAccess)
	data.Visibility = types.StringValue(organization.Visibility)
	data.Website = types.StringValue(organization.Website)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *OrganizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plannedData OrganizationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plannedData)...)
	var stateData OrganizationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}
	request := client.OrganizationPatchRequest{
		RepoAdminChangeTeamAccess: plannedData.RepoAdminChangeTeamAccess.ValueBool(),
	}
	if !plannedData.Description.IsUnknown() {
		request.Description = plannedData.Description.ValueString()
	}
	if !plannedData.Email.IsUnknown() {
		request.Email = plannedData.Email.ValueString()
	}
	if !plannedData.FullName.IsUnknown() {
		request.FullName = plannedData.FullName.ValueString()
	}
	if !plannedData.Location.IsUnknown() {
		request.Location = plannedData.Location.ValueString()
	}
	if !plannedData.Visibility.IsUnknown() {
		request.Visibility = plannedData.Visibility.ValueString()
	}
	if !plannedData.Website.IsUnknown() {
		request.Website = plannedData.Website.ValueString()
	}
	organization, err := d.client.OrganizationUpdate(
		ctx,
		plannedData.Name.ValueString(),
		&request)
	if err != nil {
		resp.Diagnostics.AddError("UpdateOrganization", fmt.Sprintf("failed to update organization: %s", err))
		return
	}
	plannedData.Description = types.StringValue(organization.Description)
	plannedData.Email = types.StringValue(organization.Email)
	plannedData.FullName = types.StringValue(organization.FullName)
	plannedData.Id = types.Int64Value(organization.Id)
	plannedData.Location = types.StringValue(organization.Location)
	plannedData.RepoAdminChangeTeamAccess = types.BoolValue(organization.RepoAdminChangeTeamAccess)
	plannedData.Visibility = types.StringValue(organization.Visibility)
	plannedData.Website = types.StringValue(organization.Website)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plannedData)...)
}
