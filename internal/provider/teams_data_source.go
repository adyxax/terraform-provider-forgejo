package provider

import (
	"context"
	"fmt"
	"slices"

	"git.adyxax.org/adyxax/terraform-provider-forgejo/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TeamsDataSource struct {
	client *client.Client
}

var _ datasource.DataSource = &TeamsDataSource{} // Ensure provider defined types fully satisfy framework interfaces
func NewTeamsDataSource() datasource.DataSource {
	return &TeamsDataSource{}
}

type TeamsDataSourceModel struct {
	Elements         []TeamDataSourceModel `tfsdk:"elements"`
	OrganizationName types.String          `tfsdk:"organization_name"`
}

type TeamDataSourceModel struct {
	CanCreateOrgRepo        types.Bool   `tfsdk:"can_create_org_repo"`
	Description             types.String `tfsdk:"description"`
	Id                      types.Int64  `tfsdk:"id"`
	IncludesAllRepositories types.Bool   `tfsdk:"includes_all_repositories"`
	Name                    types.String `tfsdk:"name"`
	// Appears unused, the TeamsList function always returns nil
	//Organization            *OrganizationDataSourceModel  `tfsdk:"organization"`
	Permission types.String            `tfsdk:"permission"`
	Units      []types.String          `tfsdk:"units"`
	UnitsMap   map[string]types.String `tfsdk:"units_map"`
}

func (d *TeamsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_teams"
}

var teamSchemaAttributes = schema.ListNestedAttribute{
	Computed:            true,
	MarkdownDescription: "The list of teams for an organization.",
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"can_create_org_repo": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether members of this team can create repositories that will belong to the organization.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "A description string.",
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The identifier of the team.",
			},
			"includes_all_repositories": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether members of this team can access all the repositories that belong to the organization.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The team's name.",
			},
			"permission": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The members' permission level on the organization.",
			},
			"units": schema.ListAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "The list of units permissions.",
			},
			"units_map": schema.MapAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "The map of units permissions and their level.",
			},
		},
	},
}

func (d *TeamsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"elements": teamSchemaAttributes,
			"organization_name": schema.StringAttribute{
				MarkdownDescription: "The name of the organization the teams are a part of.",
				Required:            true,
			},
		},
		MarkdownDescription: "Use this data source to retrieve information about existing forgejo teams belonging to an organization.",
	}
}

func (d *TeamsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client, _ = req.ProviderData.(*client.Client)
}

func populateTeamDataSourceModel(team *client.OrganizationTeam) *TeamDataSourceModel {
	return &TeamDataSourceModel{
		CanCreateOrgRepo:        types.BoolValue(team.CanCreateOrgRepo),
		Description:             types.StringValue(team.Description),
		Id:                      types.Int64Value(team.Id),
		IncludesAllRepositories: types.BoolValue(team.IncludesAllRepositories),
		Name:                    types.StringValue(team.Name),
		Permission:              types.StringValue(team.Permission),
		Units:                   make([]types.String, len(team.Units)),
		UnitsMap:                make(map[string]types.String),
	}
}

func (d *TeamsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TeamsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	teams, err := d.client.OrganizationTeamsList(ctx, data.OrganizationName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("ListTeams", fmt.Sprintf("failed to list teams: %s", err))
		return
	}
	teamsList := make([]TeamDataSourceModel, len(teams))
	for i, team := range teams {
		teamsList[i] = *populateTeamDataSourceModel(&team)
		slices.Sort(team.Units)
		for j, unit := range team.Units {
			teamsList[i].Units[j] = types.StringValue(unit)
		}
		for unit, perm := range team.UnitsMap {
			teamsList[i].UnitsMap[unit] = types.StringValue(perm)
		}
	}
	data.Elements = teamsList
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
