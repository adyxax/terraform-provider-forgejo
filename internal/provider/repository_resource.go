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
	CreatedAt     timetypes.RFC3339 `tfsdk:"created_at"`
	DefaultBranch types.String      `tfsdk:"default_branch"`
	Description   types.String      `tfsdk:"description"`
	Name          types.String      `tfsdk:"name"`
	Owner         types.String      `tfsdk:"owner"`
	Private       types.Bool        `tfsdk:"private"`
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
		resp.Diagnostics.AddError("CreateRepository", fmt.Sprintf("failed to create Repository: %s", err))
		return
	}
	data.CreatedAt = timetypes.NewRFC3339TimeValue(repository.CreatedAt)
	if data.Description.IsUnknown() {
		data.Description = types.StringValue(repository.Description)
	}
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
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: owner/repository. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("owner"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
}

func (d *RepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RepositoryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	repository, err := d.client.RepositoryGet(
		ctx,
		data.Owner.ValueString(),
		data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ReadRepository", fmt.Sprintf("failed to get Repository: %s", err))
		return
	}
	data.CreatedAt = timetypes.NewRFC3339TimeValue(repository.CreatedAt)
	data.DefaultBranch = types.StringValue(repository.DefaultBranch)
	data.Description = types.StringValue(repository.Description)
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
