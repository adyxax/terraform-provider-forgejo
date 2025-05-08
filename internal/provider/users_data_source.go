package provider

import (
	"context"
	"fmt"

	"git.adyxax.org/adyxax/terraform-provider-forgejo/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type UsersDataSource struct {
	client *client.Client
}

var _ datasource.DataSource = &UsersDataSource{} // Ensure provider defined types fully satisfy framework interfaces
func NewUsersDataSource() datasource.DataSource {
	return &UsersDataSource{}
}

type UsersDataSourceModel struct {
	Elements []UserDataSourceModel `tfsdk:"elements"`
}
type UserDataSourceModel struct {
	Active           types.Bool        `tfsdk:"active"`
	AvatarUrl        types.String      `tfsdk:"avatar_url"`
	Created          timetypes.RFC3339 `tfsdk:"created"`
	Description      types.String      `tfsdk:"description"`
	Email            types.String      `tfsdk:"email"`
	FollowerCount    types.Int64       `tfsdk:"followers_count"`
	FollowingCount   types.Int64       `tfsdk:"following_count"`
	FullName         types.String      `tfsdk:"full_name"`
	HtmlUrl          types.String      `tfsdk:"html_url"`
	Id               types.Int64       `tfsdk:"id"`
	IsAdmin          types.Bool        `tfsdk:"is_admin"`
	Language         types.String      `tfsdk:"language"`
	LastLogin        timetypes.RFC3339 `tfsdk:"last_login"`
	Location         types.String      `tfsdk:"location"`
	LoginName        types.String      `tfsdk:"login_name"`
	Login            types.String      `tfsdk:"login"`
	ProhibitLogin    types.Bool        `tfsdk:"prohibit_login"`
	Pronouns         types.String      `tfsdk:"pronouns"`
	Restricted       types.Bool        `tfsdk:"restricted"`
	SourceId         types.Int64       `tfsdk:"source_id"`
	StarredRepoCount types.Int64       `tfsdk:"starred_repos_count"`
	Visibility       types.String      `tfsdk:"visibility"`
	Website          types.String      `tfsdk:"website"`
}

func (d *UsersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

func (d *UsersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"elements": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of users.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"active": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the user is active or not.",
						},
						"avatar_url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The user's avatar URL.",
						},
						"created": schema.StringAttribute{
							Computed:            true,
							CustomType:          timetypes.RFC3339Type{},
							MarkdownDescription: "The user's creation date and time.",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "A description string.",
						},
						"email": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The user's email address.",
						},
						"followers_count": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The number of followers.",
						},
						"following_count": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The number of followings.",
						},
						"full_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The user's full name.",
						},
						"html_url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The URL to this user's Forgejo profile page.",
						},
						"id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The identifier of the user.",
						},
						"is_admin": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the user is an admin or not.",
						},
						"language": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The user's chosen language.",
						},
						"last_login": schema.StringAttribute{
							Computed:            true,
							CustomType:          timetypes.RFC3339Type{},
							MarkdownDescription: "The user's last login date and time.",
						},
						"location": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The user's advertised location.",
						},
						"login_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The user's authentication sign-in name.",
						},
						"login": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The login of the user.",
						},
						"prohibit_login": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the user is allowed to log in or not.",
						},
						"pronouns": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The user's advertised pronouns.",
						},
						"restricted": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the user is restricted or not.",
						},
						"source_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The identifier of the users authentication source.",
						},
						"starred_repos_count": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The number of repositoties starred by the user.",
						},
						"visibility": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The user's visibility option: limited, private, public.",
						},
						"website": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The user's advertised website.",
						},
					},
				},
			},
		},
		MarkdownDescription: "Use this data source to retrieve information about existing forgejo users.",
	}
}

func (d *UsersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client, _ = req.ProviderData.(*client.Client)
}

func (d *UsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UsersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	users, err := d.client.UsersList(ctx)
	if err != nil {
		resp.Diagnostics.AddError("ListUsers", fmt.Sprintf("failed to list users: %s", err))
		return
	}
	userList := make([]UserDataSourceModel, len(users))
	for i, user := range users {
		userList[i] = UserDataSourceModel{
			Active:           types.BoolValue(user.Active),
			AvatarUrl:        types.StringValue(user.AvatarUrl),
			Created:          timetypes.NewRFC3339TimeValue(user.Created),
			Description:      types.StringValue(user.Description),
			Email:            types.StringValue(user.Email),
			FollowerCount:    types.Int64Value(user.FollowerCount),
			FollowingCount:   types.Int64Value(user.FollowingCount),
			FullName:         types.StringValue(user.FullName),
			HtmlUrl:          types.StringValue(user.HtmlUrl),
			Id:               types.Int64Value(user.Id),
			IsAdmin:          types.BoolValue(user.IsAdmin),
			Language:         types.StringValue(user.Language),
			LastLogin:        timetypes.NewRFC3339TimeValue(user.LastLogin),
			Location:         types.StringValue(user.Location),
			LoginName:        types.StringValue(user.LoginName),
			Login:            types.StringValue(user.Login),
			ProhibitLogin:    types.BoolValue(user.ProhibitLogin),
			Pronouns:         types.StringValue(user.Pronouns),
			Restricted:       types.BoolValue(user.Restricted),
			SourceId:         types.Int64Value(user.SourceId),
			StarredRepoCount: types.Int64Value(user.StarredRepoCount),
			Visibility:       types.StringValue(user.Visibility),
			Website:          types.StringValue(user.Website),
		}
	}
	data.Elements = userList
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
