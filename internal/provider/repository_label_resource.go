package provider

import (
	"context"
	"fmt"

	"git.adyxax.org/adyxax/terraform-provider-forgejo/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RepositoryLabelResource struct {
	client *client.Client
}

var _ resource.Resource = &RepositoryLabelResource{} // Ensure provider defined types fully satisfy framework interfaces
func NewRepositoryLabelResource() resource.Resource {
	return &RepositoryLabelResource{}
}

type RepositoryLabelResourceModel struct {
	Color       types.String `tfsdk:"color"`
	Description types.String `tfsdk:"description"`
	Exclusive   types.Bool   `tfsdk:"exclusive"`
	Id          types.Int64  `tfsdk:"id"`
	IsArchived  types.Bool   `tfsdk:"is_archived"`
	Name        types.String `tfsdk:"name"`
	Owner       types.String `tfsdk:"owner"`
	Repository  types.String `tfsdk:"repository"`
	Url         types.String `tfsdk:"url"`
}

func (d *RepositoryLabelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_label"
}

func (d *RepositoryLabelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"color": schema.StringAttribute{
				MarkdownDescription: "The label's color in lowercase rgb format, without a leading `#`. For example `207de5`.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description string.",
				Required:            true,
			},
			"exclusive": schema.BoolAttribute{
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether the label is exclusive or not. Defaults to `false`. Name the label `scope/item` to make it mutually exclusive with other `scope/` labels.",
				Optional:            true,
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The identifier of the repository label.",
			},
			"is_archived": schema.BoolAttribute{
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether the repository label is archived or not. Defaults to `false`",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The label's name.",
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
				MarkdownDescription: "The label's repository.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Required: true,
			},
			"url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The repository label's URL.",
			},
		},
		MarkdownDescription: "Use this resource to create and manage a repository label.",
	}
}

func (d *RepositoryLabelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	d.client, _ = req.ProviderData.(*client.Client)
}

func (d *RepositoryLabelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RepositoryLabelResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	request := client.RepositoryLabelCreateRequest{
		Color:       data.Color.ValueString(),
		Description: data.Description.ValueString(),
		Exclusive:   data.Exclusive.ValueBool(),
		IsArchived:  data.IsArchived.ValueBool(),
		Name:        data.Name.ValueString(),
	}
	label, err := d.client.RepositoryLabelCreate(
		ctx,
		data.Owner.ValueString(),
		data.Repository.ValueString(),
		&request)
	if err != nil {
		resp.Diagnostics.AddError("CreateRepositoryLabel", fmt.Sprintf("failed to create repository label: %s", err))
		return
	}
	data.Id = types.Int64Value(label.Id)
	data.Url = types.StringValue(label.Url)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *RepositoryLabelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RepositoryLabelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := d.client.RepositoryLabelDelete(
		ctx,
		data.Owner.ValueString(),
		data.Repository.ValueString(),
		data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("DeleteRepositoryLabel", fmt.Sprintf("failed to delete repository label: %s", err))
		return
	}
}

func (d *RepositoryLabelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RepositoryLabelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	label, err := d.client.RepositoryLabelGet(
		ctx,
		data.Owner.ValueString(),
		data.Repository.ValueString(),
		data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("ReadRepositoryLabel", fmt.Sprintf("failed to get repository label: %s", err))
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

func (d *RepositoryLabelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plannedData RepositoryLabelResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plannedData)...)
	var stateData RepositoryLabelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}
	request := client.RepositoryLabelUpdateRequest{
		Color:       plannedData.Color.ValueString(),
		Description: plannedData.Description.ValueString(),
		Exclusive:   plannedData.Exclusive.ValueBool(),
		IsArchived:  plannedData.IsArchived.ValueBool(),
		Name:        plannedData.Name.ValueString(),
	}
	label, err := d.client.RepositoryLabelUpdate(
		ctx,
		plannedData.Owner.ValueString(),
		plannedData.Repository.ValueString(),
		stateData.Id.ValueInt64(),
		&request)
	if err != nil {
		resp.Diagnostics.AddError("UpdateRepositoryLabel", fmt.Sprintf("failed to update repository label: %s", err))
		return
	}
	plannedData.Id = types.Int64Value(label.Id)
	plannedData.Url = types.StringValue(label.Url)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plannedData)...)
}
