package provider

import (
	"context"
	"fmt"
	"net/url"

	"git.adyxax.org/adyxax/terraform-provider-forgejo/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Provider struct {
	version string
}

var _ provider.Provider = &Provider{} // Ensure provider defined types fully satisfy framework interfaces.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &Provider{
			version: version,
		}
	}
}

type ProviderModel struct {
	ApiToken types.String `tfsdk:"api_token"`
	BaseURI  types.String `tfsdk:"base_uri"`
}

func (p *Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "forgejo"
	resp.Version = p.version
}

func (p *Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				MarkdownDescription: "Forgejo's api token",
				Required:            true,
				Sensitive:           true,
			},
			"base_uri": schema.StringAttribute{
				MarkdownDescription: "Forgejo's HTTP base uri",
				Required:            true,
			},
		},
	}
}

func (p *Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	baseURI, err := url.Parse(data.BaseURI.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid base_uri",
			fmt.Sprintf(
				"failed to parse base_uri: %s",
				err))
		return
	}
	client := client.NewClient(baseURI, data.ApiToken.ValueString())

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *Provider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRepositoryActionsSecretResource,
	}
}

func (p *Provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewOrganizationsDataSource,
		NewRepositoriesDataSource,
		NewTeamsDataSource,
		NewUsersDataSource,
	}
}
