package provider

import (
	"context"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hasura/go-graphql-client"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &stackletProvider{}
)

// stackletProvider is the provider implementation.
type stackletProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// stackletProviderModel maps provider schema data to a Go type.
type stackletProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	ApiKey   types.String `tfsdk:"api_key"`
}

// New creates a new provider instance
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &stackletProvider{
			version: version,
		}
	}
}

// Metadata returns the provider type name.
func (p *stackletProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "stacklet"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *stackletProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Stacklet.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Description: "The endpoint URL of the Stacklet GraphQL API. May also be provided via STACKLET_ENDPOINT environment variable.",
				Optional:    true,
			},
			"api_key": schema.StringAttribute{
				Description: "The API key for Stacklet authentication. May also be provided via STACKLET_API_KEY environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares a Stacklet API client for data sources and resources.
func (p *stackletProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Stacklet client")

	var config stackletProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown Stacklet API Endpoint",
			"The provider cannot create the Stacklet API client as there is an unknown configuration value for the Stacklet API endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the STACKLET_ENDPOINT environment variable.",
		)
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Stacklet API Key",
			"The provider cannot create the Stacklet API client as there is an unknown configuration value for the Stacklet API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the STACKLET_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	endpoint := os.Getenv("STACKLET_ENDPOINT")
	apiKey := os.Getenv("STACKLET_API_KEY")

	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}

	if !config.ApiKey.IsNull() {
		apiKey = config.ApiKey.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing Stacklet API Endpoint",
			"The provider cannot create the Stacklet API client as there is a missing or empty value for the Stacklet API endpoint. "+
				"Set the endpoint value in the configuration or use the STACKLET_ENDPOINT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing Stacklet API Key",
			"The provider cannot create the Stacklet API client as there is a missing or empty value for the Stacklet API key. "+
				"Set the api_key value in the configuration or use the STACKLET_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create an HTTP client with the authorization header
	httpClient := &http.Client{
		Transport: &authTransport{
			apiKey: apiKey,
			base:   http.DefaultTransport,
		},
	}

	// Create a new Stacklet client using the configuration values
	client := graphql.NewClient(endpoint, httpClient)

	// Make the Stacklet client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Stacklet client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *stackletProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewRepositoryDataSource,
		NewAccountDataSource,
		NewAccountDiscoveryDataSource,
		NewSSOGroupDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *stackletProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRepositoryResource,
		NewAccountResource,
		NewAccountDiscoveryResource,
		NewSSOGroupResource,
	}
}

// authTransport is an http.RoundTripper that adds authentication headers
type authTransport struct {
	apiKey string
	base   http.RoundTripper
}

// RoundTrip implements the http.RoundTripper interface
func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", t.apiKey)
	return t.base.RoundTrip(req)
}
