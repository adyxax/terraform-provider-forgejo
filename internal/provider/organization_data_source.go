package provider

import (
	"context"
	"fmt"

	"git.adyxax.org/adyxax/terraform-provider-forgejo/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type OrganizationDataSource struct {
	client *client.Client
}

var _ datasource.DataSource = &OrganizationDataSource{} // Ensure provider defined types fully satisfy framework interfaces
func NewOrganizationDataSource() datasource.DataSource {
	return &OrganizationDataSource{}
}

type OrganizationDataSourceModel struct {
	AvatarUrl                 types.String `tfsdk:"avatar_url"`
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

func (d *OrganizationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (d *OrganizationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"avatar_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The organization's avatar URL.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "A description string.",
			},
			"email": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The organization's email address.",
			},
			"full_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The organization's full name.",
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The identifier of the organization.",
			},
			"location": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The organization's advertised location.",
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the organization.",
				Required:            true,
			},
			"repo_admin_change_team_access": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether an admin of a repository that belongs to this organization can change team access or not.",
			},
			"visibility": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The organization's visibility option: limited, private, public.",
			},
			"website": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The organization's advertised website.",
			},
		},
		MarkdownDescription: "Use this data source to retrieve information about existing forgejo organizations.",
	}
}

func (d *OrganizationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client, _ = req.ProviderData.(*client.Client)
}

func (d *OrganizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	organization, err := d.client.OrganizationGet(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("OrganizationGet", fmt.Sprintf("failed to get organization: %s", err))
		return
	}
	data.AvatarUrl = types.StringValue(organization.AvatarUrl)
	data.Description = types.StringValue(organization.Description)
	data.Email = types.StringValue(organization.Email)
	data.Id = types.Int64Value(organization.Id)
	data.Location = types.StringValue(organization.Location)
	data.Name = types.StringValue(organization.Name)
	data.RepoAdminChangeTeamAccess = types.BoolValue(organization.RepoAdminChangeTeamAccess)
	data.Visibility = types.StringValue(organization.Visibility)
	data.Website = types.StringValue(organization.Website)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
