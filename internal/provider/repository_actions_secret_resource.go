package provider

import (
	"context"
	"fmt"
	"strings"

	"git.adyxax.org/adyxax/terraform-provider-forgejo/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RepositoryActionsSecretResource struct {
	client *client.Client
}

var _ resource.Resource = &RepositoryActionsSecretResource{} // Ensure provider defined types fully satisfy framework interfaces
func NewRepositoryActionsSecretResource() resource.Resource {
	return &RepositoryActionsSecretResource{}
}

type RepositoryActionsSecretResourceModel struct {
	CreatedAt  timetypes.RFC3339 `tfsdk:"created_at"`
	Data       types.String      `tfsdk:"data"`
	Name       types.String      `tfsdk:"name"`
	Owner      types.String      `tfsdk:"owner"`
	Repository types.String      `tfsdk:"repository"`
}

func (d *RepositoryActionsSecretResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_actions_secret"
}

func (d *RepositoryActionsSecretResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"created_at": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The secret's creation date and time.",
			},
			"data": schema.StringAttribute{
				MarkdownDescription: "The secret's data.",
				Required:            true,
				Sensitive:           true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The secret's name. It must be uppercase or the plan will not be idempotent.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Required: true,
			},
			"owner": schema.StringAttribute{
				MarkdownDescription: "The secret's owner.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Required: true,
			},
			"repository": schema.StringAttribute{
				MarkdownDescription: "The secret's repository.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Required: true,
			},
		},
		MarkdownDescription: "Use this resource to create and manage a repository actions secret.",
	}
}

func (d *RepositoryActionsSecretResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	d.client, _ = req.ProviderData.(*client.Client)
}

func (d *RepositoryActionsSecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RepositoryActionsSecretResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := d.client.RepositoryActionsSecretCreateOrUpdate(
		ctx,
		data.Owner.ValueString(),
		data.Repository.ValueString(),
		data.Name.ValueString(),
		data.Data.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("CreateRepositoryActionsSecret", fmt.Sprintf("failed to create or update repository actions secret: %s", err))
		return
	}
	secret, err := d.getRepositoryActionsSecret(ctx, data.Owner, data.Repository, data.Name)
	if err != nil {
		resp.Diagnostics.AddError("CreateRepositoryActionsSecret", err.Error())
		return
	}
	data.CreatedAt = secret.CreatedAt
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *RepositoryActionsSecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RepositoryActionsSecretResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := d.client.RepositoryActionsSecretDelete(
		ctx,
		data.Owner.ValueString(),
		data.Repository.ValueString(),
		data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("DeleteRepositoryActionsSecret", fmt.Sprintf("failed to delete repository actions secret: %s", err))
		return
	}
}

func (d *RepositoryActionsSecretResource) getRepositoryActionsSecret(
	ctx context.Context,
	owner types.String,
	repository types.String,
	name types.String,
) (*RepositoryActionsSecretResourceModel, error) {
	secrets, err := d.client.RepositoryActionsSecretsList(
		ctx,
		owner.ValueString(),
		repository.ValueString())
	if err != nil {
		return nil, fmt.Errorf("failed to list repository actions secrets: %w", err)
	}
	nameStr := strings.ToUpper(name.ValueString())
	for _, secret := range secrets {
		if secret.Name == nameStr {
			created := timetypes.NewRFC3339TimeValue(secret.CreatedAt)
			return &RepositoryActionsSecretResourceModel{
				CreatedAt: created,
				Name:      types.StringValue(secret.Name),
			}, nil
		}
	}
	return nil, fmt.Errorf("failed to find repository actions secret")
}

func (d *RepositoryActionsSecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RepositoryActionsSecretResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	secret, err := d.getRepositoryActionsSecret(ctx, data.Owner, data.Repository, data.Name)
	if err != nil {
		resp.Diagnostics.AddError("CreateRepositoryActionsSecret", err.Error())
		return
	}
	data.CreatedAt = secret.CreatedAt
	data.Name = secret.Name
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *RepositoryActionsSecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RepositoryActionsSecretResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := d.client.RepositoryActionsSecretCreateOrUpdate(
		ctx,
		data.Owner.ValueString(),
		data.Repository.ValueString(),
		data.Name.ValueString(),
		data.Data.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("UpdateRepositoryActionsSecret", fmt.Sprintf("failed to create or update repository actions secret: %s", err))
		return
	}
	secret, err := d.getRepositoryActionsSecret(ctx, data.Owner, data.Repository, data.Name)
	if err != nil {
		resp.Diagnostics.AddError("UpdateRepositoryActionsSecret", err.Error())
		return
	}
	data.CreatedAt = secret.CreatedAt
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
