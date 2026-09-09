package provider

import (
	"context"
	"fmt"

	"git.adyxax.org/adyxax/terraform-provider-forgejo/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
				MarkdownDescription: "A description string. If unset, the server value will be left as is.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The organization's email address. If unset, the server value will be left as is.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"full_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The organization's full name. If unset, the server value will be left as is.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The identifier of the organization.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"location": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The organization's advertised location. If unset, the server value will be left as is.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				MarkdownDescription: "The organization's visibility option: `limited`, `private`, `public`. Defaults to `private`. Because of a forgejo API limitation, switching an existing organization to `public` requires administrator credentials.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("limited", "private", "public"),
				},
			},
			"website": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The organization's advertised website. If unset, the server value will be left as is.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
		Description:               stringPointerIfKnown(data.Description),
		Email:                     stringPointerIfKnown(data.Email),
		FullName:                  stringPointerIfKnown(data.FullName),
		Location:                  stringPointerIfKnown(data.Location),
		RepoAdminChangeTeamAccess: boolPointerIfKnown(data.RepoAdminChangeTeamAccess),
		Username:                  data.Name.ValueString(),
		Visibility:                data.Visibility.ValueString(),
		Website:                   stringPointerIfKnown(data.Website),
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
		Description:               stringPointerIfKnown(plannedData.Description),
		Email:                     stringPointerIfKnown(plannedData.Email),
		FullName:                  stringPointerIfKnown(plannedData.FullName),
		Location:                  stringPointerIfKnown(plannedData.Location),
		RepoAdminChangeTeamAccess: boolPointerIfKnown(plannedData.RepoAdminChangeTeamAccess),
		Visibility:                plannedData.Visibility.ValueString(),
		Website:                   stringPointerIfKnown(plannedData.Website),
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
