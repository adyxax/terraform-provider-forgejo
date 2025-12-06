package provider

import (
	"context"
	"fmt"

	"git.adyxax.org/adyxax/terraform-provider-forgejo/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type OrganizationsDataSource struct {
	client *client.Client
}

var _ datasource.DataSource = &OrganizationsDataSource{} // Ensure provider defined types fully satisfy framework interfaces
func NewOrganizationsDataSource() datasource.DataSource {
	return &OrganizationsDataSource{}
}

type OrganizationsDataSourceModel struct {
	Elements []OrganizationDataSourceModel `tfsdk:"elements"`
}

func (d *OrganizationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations"
}

func (d *OrganizationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"elements": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of organizations.",
				NestedObject: schema.NestedAttributeObject{
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
							Computed:            true,
							MarkdownDescription: "The name of the organization.",
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
				},
			},
		},
		MarkdownDescription: "Use this data source to retrieve information about existing forgejo organizations.",
	}
}

func (d *OrganizationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client, _ = req.ProviderData.(*client.Client)
}

func (d *OrganizationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	organizations, err := d.client.OrganizationsList(ctx)
	if err != nil {
		resp.Diagnostics.AddError("OrganizationsList", fmt.Sprintf("failed to list organizations: %s", err))
		return
	}
	organizationList := make([]OrganizationDataSourceModel, len(organizations))
	for i, organization := range organizations {
		organizationList[i] = OrganizationDataSourceModel{
			AvatarUrl:                 types.StringValue(organization.AvatarUrl),
			Description:               types.StringValue(organization.Description),
			Email:                     types.StringValue(organization.Email),
			Id:                        types.Int64Value(organization.Id),
			Location:                  types.StringValue(organization.Location),
			Name:                      types.StringValue(organization.Name),
			RepoAdminChangeTeamAccess: types.BoolValue(organization.RepoAdminChangeTeamAccess),
			Visibility:                types.StringValue(organization.Visibility),
			Website:                   types.StringValue(organization.Website),
		}
	}
	data.Elements = organizationList
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
