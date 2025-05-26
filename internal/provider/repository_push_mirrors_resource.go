package provider

import (
	"context"
	"fmt"

	"git.adyxax.org/adyxax/terraform-provider-forgejo/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RepositoryPushMirrorResource struct {
	client *client.Client
}

var _ resource.Resource = &RepositoryPushMirrorResource{} // Ensure provider defined types fully satisfy framework interfaces
func NewRepositoryPushMirrorResource() resource.Resource {
	return &RepositoryPushMirrorResource{}
}

type RepositoryPushMirrorResourceModel struct {
	Created        timetypes.RFC3339 `tfsdk:"created"`
	Interval       types.String      `tfsdk:"interval"`
	Name           types.String      `tfsdk:"name"`
	Owner          types.String      `tfsdk:"owner"`
	RemoteAddress  types.String      `tfsdk:"remote_address"`
	RemotePassword types.String      `tfsdk:"remote_password"`
	RemoteUsername types.String      `tfsdk:"remote_username"`
	Repository     types.String      `tfsdk:"repository"`
	SyncOnCommit   types.Bool        `tfsdk:"sync_on_commit"`
	UseSsh         types.Bool        `tfsdk:"use_ssh"`
}

func (d *RepositoryPushMirrorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_push_mirror"
}

func (d *RepositoryPushMirrorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"created": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The push mirror's creation date and time.",
			},
			"interval": schema.StringAttribute{
				Computed:            true,
				Default:             stringdefault.StaticString("8h0m0s"),
				MarkdownDescription: "The push mirror's sync interval as a string. Defaults to `8h0m0s`.",
				Optional:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the push mirror.",
			},
			"owner": schema.StringAttribute{
				MarkdownDescription: "The owner of the repository on which to configure a push mirror.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Required:            true,
			},
			"remote_address": schema.StringAttribute{
				MarkdownDescription: "The push mirror's remote address.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Required:            true,
			},
			"remote_password": schema.StringAttribute{
				MarkdownDescription: "The push mirror's remote password.",
				Optional:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Sensitive:           true,
			},
			"remote_username": schema.StringAttribute{
				MarkdownDescription: "The push mirror's remote username.",
				Optional:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"repository": schema.StringAttribute{
				MarkdownDescription: "The repository on which to configure a push mirror.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Required:            true,
			},
			"sync_on_commit": schema.BoolAttribute{
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				MarkdownDescription: "Whether the push mirror is synced on each commit pushed to the repository, defaults to `true`.",
				Optional:            true,
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
			"use_ssh": schema.BoolAttribute{
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Whether the push mirror is synced over SSH or not (not meaning HTTP), defaults to `false`.",
				Optional:            true,
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
		},
		MarkdownDescription: "Use this resource to create and manage a repository push mirror.",
	}
}

func (d *RepositoryPushMirrorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	d.client, _ = req.ProviderData.(*client.Client)
}

func (d *RepositoryPushMirrorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RepositoryPushMirrorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	pushMirror, err := d.client.RepositoryPushMirrorCreate(
		ctx,
		data.Owner.ValueString(),
		data.Repository.ValueString(),
		data.Interval.ValueString(),
		data.RemoteAddress.ValueString(),
		data.RemotePassword.ValueString(),
		data.RemoteUsername.ValueString(),
		data.SyncOnCommit.ValueBool(),
		data.UseSsh.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("CreateRepositoryPushMirror", fmt.Sprintf("failed to create repository push mirror: %s", err))
		return
	}
	data.Created = timetypes.NewRFC3339TimeValue(pushMirror.Created)
	data.Name = types.StringValue(pushMirror.RemoteName)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *RepositoryPushMirrorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RepositoryPushMirrorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, err := d.client.RepositoryPushMirrorDelete(
		ctx,
		data.Owner.ValueString(),
		data.Repository.ValueString(),
		data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("DeleteRepositoryPushMirror", fmt.Sprintf("failed to delete repository push mirror: %s", err))
		return
	}
}

func (d *RepositoryPushMirrorResource) getRepositoryPushMirror(
	ctx context.Context,
	owner types.String,
	repository types.String,
	name types.String,
) (*client.RepositoryPushMirror, error) {
	pushMirror, err := d.client.RepositoryPushMirrorGet(
		ctx,
		owner.ValueString(),
		repository.ValueString(),
		name.ValueString(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository push mirror: %w", err)
	}
	return pushMirror, nil
}

func (d *RepositoryPushMirrorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RepositoryPushMirrorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	pushMirror, err := d.getRepositoryPushMirror(ctx, data.Owner, data.Repository, data.Name)
	if err != nil {
		resp.Diagnostics.AddError("ReadRepositoryPushMirror", fmt.Sprintf("failed to get repository push mirror: %s", err))
		return
	}
	data.Created = timetypes.NewRFC3339TimeValue(pushMirror.Created)
	data.Name = types.StringValue(pushMirror.RemoteName)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *RepositoryPushMirrorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("UpdateRepositoryPushMirror", "unreachable code")
}
