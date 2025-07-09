package provider

import (
	"context"
	"fmt"
	"strings"

	"git.adyxax.org/adyxax/terraform-provider-forgejo/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RepositoryActionsVariableResource struct {
	client *client.Client
}

var _ resource.Resource = &RepositoryActionsVariableResource{}                // Ensure provider defined types fully satisfy framework interfaces
var _ resource.ResourceWithImportState = &RepositoryActionsVariableResource{} // Ensure provider defined types fully satisfy framework interfaces
func NewRepositoryActionsVariableResource() resource.Resource {
	return &RepositoryActionsVariableResource{}
}

type RepositoryActionsVariableResourceModel struct {
	Data       types.String `tfsdk:"data"`
	Name       types.String `tfsdk:"name"`
	Owner      types.String `tfsdk:"owner"`
	Repository types.String `tfsdk:"repository"`
}

func (d *RepositoryActionsVariableResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_actions_variable"
}

func (d *RepositoryActionsVariableResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"data": schema.StringAttribute{
				MarkdownDescription: "The variable's data.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The variable's name. It must be uppercase or the plan will not be idempotent.",
				Required:            true,
			},
			"owner": schema.StringAttribute{
				MarkdownDescription: "The variable's owner.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Required: true,
			},
			"repository": schema.StringAttribute{
				MarkdownDescription: "The variable's repository.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Required: true,
			},
		},
		MarkdownDescription: "Use this resource to create and manage a repository actions variable.",
	}
}

func (d *RepositoryActionsVariableResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	d.client, _ = req.ProviderData.(*client.Client)
}

func (d *RepositoryActionsVariableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RepositoryActionsVariableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := d.client.RepositoryActionsVariableCreate(
		ctx,
		data.Owner.ValueString(),
		data.Repository.ValueString(),
		data.Name.ValueString(),
		data.Data.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("CreateRepositoryActionsVariable", fmt.Sprintf("failed to create repository actions variable: %s\nTry importing the resource instead?", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *RepositoryActionsVariableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RepositoryActionsVariableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, err := d.client.RepositoryActionsVariableDelete(
		ctx,
		data.Owner.ValueString(),
		data.Repository.ValueString(),
		data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("DeleteRepositoryActionsVariable", fmt.Sprintf("failed to delete repository actions variable: %s", err))
		return
	}
}

func (r *RepositoryActionsVariableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "/")
	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: owner/repository/variableName. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("owner"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("repository"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[2])...)
}

func (d *RepositoryActionsVariableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RepositoryActionsVariableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	variable, err := d.client.RepositoryActionsVariableGet(
		ctx,
		data.Owner.ValueString(),
		data.Repository.ValueString(),
		data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ReadRepositoryActionsVariable", fmt.Sprintf("failed to get repository actions variable: %s", err))
		return
	}
	data.Data = types.StringValue(variable.Data)
	data.Name = types.StringValue(variable.Name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *RepositoryActionsVariableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plannedData RepositoryActionsVariableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plannedData)...)
	var stateData RepositoryActionsVariableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := d.client.RepositoryActionsVariableUpdate(
		ctx,
		plannedData.Owner.ValueString(),
		plannedData.Repository.ValueString(),
		stateData.Name.ValueString(),
		plannedData.Name.ValueString(),
		plannedData.Data.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("UpdateRepositoryActionsVariable", fmt.Sprintf("failed to update repository actions variable: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plannedData)...)
}
