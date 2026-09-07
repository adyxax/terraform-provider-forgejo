package provider

import (
	"context"
	"fmt"

	"git.adyxax.org/adyxax/terraform-provider-forgejo/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type UserResource struct {
	client *client.Client
}

var _ resource.Resource = &UserResource{}                // Ensure provider defined types fully satisfy framework interfaces
var _ resource.ResourceWithImportState = &UserResource{} // Ensure provider defined types fully satisfy framework interfaces
func NewUserResource() resource.Resource {
	return &UserResource{}
}

type UserResourceModel struct {
	Active             types.Bool        `tfsdk:"active"`
	Created            timetypes.RFC3339 `tfsdk:"created"`
	Description        types.String      `tfsdk:"description"`
	Email              types.String      `tfsdk:"email"`
	FullName           types.String      `tfsdk:"full_name"`
	HtmlUrl            types.String      `tfsdk:"html_url"`
	Id                 types.Int64       `tfsdk:"id"`
	IsAdmin            types.Bool        `tfsdk:"is_admin"`
	Location           types.String      `tfsdk:"location"`
	LoginName          types.String      `tfsdk:"login_name"`
	Login              types.String      `tfsdk:"login"`
	MustChangePassword types.Bool        `tfsdk:"must_change_password"`
	Password           types.String      `tfsdk:"password"`
	ProhibitLogin      types.Bool        `tfsdk:"prohibit_login"`
	Pronouns           types.String      `tfsdk:"pronouns"`
	Restricted         types.Bool        `tfsdk:"restricted"`
	SourceId           types.Int64       `tfsdk:"source_id"`
	Visibility         types.String      `tfsdk:"visibility"`
	Website            types.String      `tfsdk:"website"`
}

func (d *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"active": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the user is active or not. If unset, the server value will be left as is.",
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"created": schema.StringAttribute{
				Computed:            true,
				CustomType:          timetypes.RFC3339Type{},
				MarkdownDescription: "The user's creation date and time.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "A description string. If unset, the server value will be left as is.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The user's email address.",
				Required:            true,
			},
			"full_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The user's full name. If unset, the server value will be left as is.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"html_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The URL to this user's Forgejo profile page.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The identifier of the user.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"is_admin": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the user is an admin or not. If unset, the server value will be left as is.",
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"location": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The user's advertised location. If unset, the server value will be left as is.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"login_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The user's authentication sign-in name, used with external authentication sources. Defaults to the user's login.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"login": schema.StringAttribute{
				MarkdownDescription: "The login (username) of the user.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Required: true,
			},
			"must_change_password": schema.BoolAttribute{
				MarkdownDescription: "Whether the user must change their password on next sign-in. If unset, the server value will be left as is. This value cannot be read back from forgejo: changes made outside of terraform will not be detected.",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The user's password. Required when `source_id` is `0` (local authentication). This value cannot be read back from forgejo: changes made outside of terraform will not be detected.",
				Optional:            true,
				Sensitive:           true,
			},
			"prohibit_login": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the user is prohibited from signing in or not. If unset, the server value will be left as is.",
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"pronouns": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The user's advertised pronouns. If unset, the server value will be left as is.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"restricted": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the user is restricted or not. If unset, the server value will be left as is.",
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"source_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The identifier of the user's authentication source, `0` being local authentication. If unset, the server value will be left as is.",
				Optional:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"visibility": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The user's visibility option: `limited`, `private`, `public`. If unset, the server value will be left as is.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("limited", "private", "public"),
				},
			},
			"website": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The user's advertised website. If unset, the server value will be left as is.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		MarkdownDescription: "Use this resource to create and manage a user. The provider must be configured with the credentials of an administrator to use this resource.",
	}
}

func (d *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	d.client, _ = req.ProviderData.(*client.Client)
}

func populateUserResourceModel(data *UserResourceModel, user *client.User) {
	data.Active = types.BoolValue(user.Active)
	data.Created = timetypes.NewRFC3339TimeValue(user.Created)
	data.Description = types.StringValue(user.Description)
	data.Email = types.StringValue(user.Email)
	data.FullName = types.StringValue(user.FullName)
	data.HtmlUrl = types.StringValue(user.HtmlUrl)
	data.Id = types.Int64Value(user.Id)
	data.IsAdmin = types.BoolValue(user.IsAdmin)
	data.Location = types.StringValue(user.Location)
	data.LoginName = types.StringValue(user.LoginName)
	data.Login = types.StringValue(user.Login)
	// data.Password is ignored because the user GET API call obviously does not return it
	data.ProhibitLogin = types.BoolValue(user.ProhibitLogin)
	data.Pronouns = types.StringValue(user.Pronouns)
	data.Restricted = types.BoolValue(user.Restricted)
	data.SourceId = types.Int64Value(user.SourceId)
	data.Visibility = types.StringValue(user.Visibility)
	data.Website = types.StringValue(user.Website)
}

func userUpdateRequest(data *UserResourceModel, sendPassword bool) *client.UserUpdateRequest {
	loginName := data.LoginName.ValueString()
	if loginName == "" {
		loginName = data.Login.ValueString()
	}
	request := client.UserUpdateRequest{
		Active:             boolPointerIfKnown(data.Active),
		Admin:              boolPointerIfKnown(data.IsAdmin),
		Description:        stringPointerIfKnown(data.Description),
		Email:              data.Email.ValueStringPointer(),
		FullName:           stringPointerIfKnown(data.FullName),
		Location:           stringPointerIfKnown(data.Location),
		LoginName:          loginName,
		MustChangePassword: boolPointerIfKnown(data.MustChangePassword),
		ProhibitLogin:      boolPointerIfKnown(data.ProhibitLogin),
		Pronouns:           stringPointerIfKnown(data.Pronouns),
		Restricted:         boolPointerIfKnown(data.Restricted),
		SourceId:           data.SourceId.ValueInt64(),
		Visibility:         data.Visibility.ValueString(),
		Website:            stringPointerIfKnown(data.Website),
	}
	if sendPassword && !data.Password.IsNull() && !data.Password.IsUnknown() {
		request.Password = data.Password.ValueString()
	}
	return &request
}

func (d *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if data.SourceId.ValueInt64() == 0 && data.Password.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("password"), "CreateUser", "a password is required when source_id is 0 (local authentication)")
		return
	}
	createRequest := client.UserCreateRequest{
		Email:              data.Email.ValueString(),
		FullName:           data.FullName.ValueString(),
		MustChangePassword: boolPointerIfKnown(data.MustChangePassword),
		Password:           data.Password.ValueString(),
		Restricted:         boolPointerIfKnown(data.Restricted),
		SendNotify:         false,
		SourceId:           data.SourceId.ValueInt64(),
		Username:           data.Login.ValueString(),
		Visibility:         data.Visibility.ValueString(),
	}
	if !data.LoginName.IsUnknown() && !data.LoginName.IsNull() {
		createRequest.LoginName = data.LoginName.ValueString()
	}
	user, err := d.client.UserCreate(ctx, &createRequest)
	if err != nil {
		resp.Diagnostics.AddError("CreateUser", fmt.Sprintf("failed to create user: %s", err))
		return
	}
	if data.LoginName.IsUnknown() || data.LoginName.IsNull() {
		data.LoginName = types.StringValue(user.LoginName)
	}
	user, err = d.client.UserUpdate(ctx, user.Login, userUpdateRequest(&data, false))
	if err != nil {
		resp.Diagnostics.AddError("CreateUser", fmt.Sprintf("failed to update user after creation: %s", err))
		return
	}
	populateUserResourceModel(&data, user)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := d.client.UserDelete(
		ctx,
		data.Login.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("DeleteUser", fmt.Sprintf("failed to delete user: %s", err))
		return
	}
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError("ImportUser", "Expected import identifier to be the user's login")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("login"), req.ID)...)
}

func (d *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	user, err := d.client.UserGet(
		ctx,
		data.Login.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ReadUser", fmt.Sprintf("failed to get user: %s", err))
		return
	}
	populateUserResourceModel(&data, user)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plannedData UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plannedData)...)
	var stateData UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}
	passwordChanged := !plannedData.Password.Equal(stateData.Password)
	if plannedData.LoginName.IsUnknown() || plannedData.LoginName.IsNull() {
		plannedData.LoginName = stateData.LoginName
	}
	user, err := d.client.UserUpdate(
		ctx,
		stateData.Login.ValueString(),
		userUpdateRequest(&plannedData, passwordChanged))
	if err != nil {
		resp.Diagnostics.AddError("UpdateUser", fmt.Sprintf("failed to update user: %s", err))
		return
	}
	var serverData UserResourceModel
	populateUserResourceModel(&serverData, user)
	plannedData.Active = knownOr(plannedData.Active, serverData.Active)
	plannedData.Created = knownOr(plannedData.Created, serverData.Created)
	plannedData.Description = knownOr(plannedData.Description, serverData.Description)
	plannedData.Email = knownOr(plannedData.Email, serverData.Email)
	plannedData.FullName = knownOr(plannedData.FullName, serverData.FullName)
	plannedData.HtmlUrl = knownOr(plannedData.HtmlUrl, serverData.HtmlUrl)
	plannedData.Id = knownOr(plannedData.Id, serverData.Id)
	plannedData.IsAdmin = knownOr(plannedData.IsAdmin, serverData.IsAdmin)
	plannedData.Location = knownOr(plannedData.Location, serverData.Location)
	plannedData.LoginName = knownOr(plannedData.LoginName, serverData.LoginName)
	plannedData.Login = knownOr(plannedData.Login, serverData.Login)
	plannedData.ProhibitLogin = knownOr(plannedData.ProhibitLogin, serverData.ProhibitLogin)
	plannedData.Pronouns = knownOr(plannedData.Pronouns, serverData.Pronouns)
	plannedData.Restricted = knownOr(plannedData.Restricted, serverData.Restricted)
	plannedData.SourceId = knownOr(plannedData.SourceId, serverData.SourceId)
	plannedData.Visibility = knownOr(plannedData.Visibility, serverData.Visibility)
	plannedData.Website = knownOr(plannedData.Website, serverData.Website)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plannedData)...)
}
