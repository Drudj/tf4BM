package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/selectel/terraform-provider-selectel-baremetal/internal/client"
	"github.com/selectel/terraform-provider-selectel-baremetal/internal/datasources"
	"github.com/selectel/terraform-provider-selectel-baremetal/internal/resources"
)

// Ensure SelectelBaremetalProvider satisfies various provider interfaces.
var _ provider.Provider = &SelectelBaremetalProvider{}

// SelectelBaremetalProvider defines the provider implementation.
type SelectelBaremetalProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// SelectelBaremetalProviderModel describes the provider data model.
type SelectelBaremetalProviderModel struct {
	Token     types.String `tfsdk:"token"`
	ProjectID types.String `tfsdk:"project_id"`
	Endpoint  types.String `tfsdk:"endpoint"`
}

func (p *SelectelBaremetalProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "selectel-baremetal"
	resp.Version = p.version
}

func (p *SelectelBaremetalProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Terraform provider для управления выделенными серверами (bare metal) Selectel.",
		Description:         "Terraform provider для управления выделенными серверами (bare metal) Selectel.",
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				MarkdownDescription: "API токен Selectel для аутентификации. Может быть задан через переменную окружения `SELECTEL_TOKEN`.",
				Description:         "API токен Selectel для аутентификации.",
				Optional:            true,
				Sensitive:           true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "UUID проекта Selectel. Может быть задан через переменную окружения `SELECTEL_PROJECT_ID`.",
				Description:         "UUID проекта Selectel.",
				Optional:            true,
			},
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "URL эндпоинта API выделенных серверов Selectel. По умолчанию: `https://api.selectel.ru/dedicated/v2`.",
				Description:         "URL эндпоинта API выделенных серверов Selectel.",
				Optional:            true,
			},
		},
	}
}

func (p *SelectelBaremetalProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data SelectelBaremetalProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// Example code for setting up client configuration

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if data.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Selectel API Token",
			"The provider cannot create the Selectel API client as there is an unknown configuration value for the Selectel API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SELECTEL_TOKEN environment variable.",
		)
	}

	if data.ProjectID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("project_id"),
			"Unknown Selectel Project ID",
			"The provider cannot create the Selectel API client as there is an unknown configuration value for the Selectel project ID. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SELECTEL_PROJECT_ID environment variable.",
		)
	}

	if data.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown Selectel API Endpoint",
			"The provider cannot create the Selectel API client as there is an unknown configuration value for the Selectel API endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SELECTEL_ENDPOINT environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	token := os.Getenv("SELECTEL_TOKEN")
	projectID := os.Getenv("SELECTEL_PROJECT_ID")
	endpoint := os.Getenv("SELECTEL_ENDPOINT")

	if !data.Token.IsNull() {
		token = data.Token.ValueString()
	}

	if !data.ProjectID.IsNull() {
		projectID = data.ProjectID.ValueString()
	}

	if !data.Endpoint.IsNull() {
		endpoint = data.Endpoint.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Selectel API Token",
			"The provider requires a Selectel API token to authenticate with the Selectel API. "+
				"Set the token value in the configuration or use the SELECTEL_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if projectID == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("project_id"),
			"Missing Selectel Project ID",
			"The provider requires a Selectel project ID to manage resources. "+
				"Set the project_id value in the configuration or use the SELECTEL_PROJECT_ID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if endpoint == "" {
		endpoint = "https://api.selectel.ru/dedicated/v2"
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "selectel_token", token)
	ctx = tflog.SetField(ctx, "selectel_project_id", projectID)
	ctx = tflog.SetField(ctx, "selectel_endpoint", endpoint)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "selectel_token")

	tflog.Debug(ctx, "Creating Selectel client")

	// Create API client
	apiClient := client.NewAPIClient(client.Config{
		Token:     token,
		ProjectID: projectID,
		Endpoint:  endpoint,
	})

	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient

	tflog.Info(ctx, "Configured Selectel client", map[string]any{"success": true})
}

func (p *SelectelBaremetalProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewServerResource,
	}
}

func (p *SelectelBaremetalProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewLocationsDataSource,
		datasources.NewLocationDataSource,
		datasources.NewServicesDataSource,
		datasources.NewServiceDataSource,
		datasources.NewOSTemplatesDataSource,
		datasources.NewOSTemplateDataSource,
		datasources.NewPricePlansDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SelectelBaremetalProvider{
			version: version,
		}
	}
}
