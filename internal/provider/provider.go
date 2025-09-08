package provider

import (
	"context"
	"fmt"
	"net/url"
	"os"

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
				MarkdownDescription: "Forgejo's api token. If not defined, the content of the environment variable `FORGEJO_API_TOKEN` will be used instead.",
				Optional:            true,
				Sensitive:           true,
			},
			"base_uri": schema.StringAttribute{
				MarkdownDescription: "Forgejo's HTTP base URI.",
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
		resp.Diagnostics.AddError("Invalid forgejo base_uri", fmt.Sprintf("failed to parse base_uri: %s", err))
		return
	}
	var apiToken string
	if data.ApiToken.IsNull() {
		apiToken = os.Getenv("FORGEJO_API_TOKEN")
		if apiToken == "" {
			resp.Diagnostics.AddError("Invalid forgejo api_token", "environment variable FORGEJO_API_TOKEN not found")
			return
		}
	} else {
		apiToken = data.ApiToken.ValueString()
	}
	client, err := client.NewClient(ctx, baseURI, apiToken)
	if err != nil {
		resp.Diagnostics.AddError("failed to instantiate forgejo client", err.Error())
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *Provider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRepositoryActionsSecretResource,
		NewRepositoryActionsVariableResource,
		NewRepositoryLabelResource,
		NewRepositoryPushMirrorResource,
		NewRepositoryResource,
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
