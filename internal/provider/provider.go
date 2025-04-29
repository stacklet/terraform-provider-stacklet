package provider

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"path"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	tfpath "github.com/hashicorp/terraform-plugin-framework/path"
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
			tfpath.Root("endpoint"),
			"Unknown Stacklet API Endpoint",
			"The provider cannot create the Stacklet API client as there is an unknown configuration value for the Stacklet API endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the STACKLET_ENDPOINT environment variable.",
		)
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			tfpath.Root("api_key"),
			"Unknown Stacklet API Key",
			"The provider cannot create the Stacklet API client as there is an unknown configuration value for the Stacklet API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the STACKLET_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	creds := getCredentials(config)

	if creds.endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			tfpath.Root("endpoint"),
			"Missing Stacklet API Endpoint",
			"The provider cannot create the Stacklet API client as there is a missing or empty value for the Stacklet API endpoint. "+
				"Set the endpoint value in the configuration, in the STACKLET_ENDPOINT environment variable, or login via the stacklet-admin CLI first.",
		)
	}
	if creds.apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			tfpath.Root("api_key"),
			"Missing Stacklet API Key",
			"The provider cannot create the Stacklet API client as there is a missing or empty value for the Stacklet API key. "+
				"Set the api_key value in the configuration, in the STACKLET_API_KEY environment variable, or login via the stacklet-admin CLI first.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create an HTTP client with the authorization header
	httpClient := &http.Client{
		Transport: &authTransport{
			apiKey: creds.apiKey,
			base:   http.DefaultTransport,
		},
	}

	// Create a new Stacklet client using the configuration values
	client := graphql.NewClient(creds.endpoint, httpClient)

	// Make the Stacklet client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Stacklet client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *stackletProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAccountDataSource,
		NewAccountGroupDataSource,
		NewPolicyDataSource,
		NewPolicyCollectionDataSource,
		NewBindingDataSource,
		NewAccountDiscoveryDataSource,
		NewSSOGroupDataSource,
		NewRepositoryDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *stackletProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAccountResource,
		NewAccountGroupResource,
		NewAccountGroupItemResource,
		NewPolicyCollectionItemResource,
		NewPolicyCollectionResource,
		NewBindingResource,
		NewAccountDiscoveryResource,
		NewSSOGroupResource,
		NewRepositoryResource,
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

func wrapNodeID(parts []string) graphql.ID {
	jsonBytes, err := json.Marshal(parts)
	if err != nil {
		// This should never happen with a simple string array
		return graphql.ID("")
	}
	encoded := base64.StdEncoding.EncodeToString(jsonBytes)
	return graphql.ID(encoded)
}

type credentials struct {
	endpoint string
	apiKey   string
}

// getCredentials tries to obtain provider credentials. The following sources
// are looked up in order:
// - provider configuration parameters
// - environment variables
// - stacklet-admin configuration
func getCredentials(config stackletProviderModel) *credentials {
	creds := credentials{}

	// Lookup provider configuration
	if !config.Endpoint.IsNull() {
		creds.endpoint = config.Endpoint.ValueString()
	}
	if !config.ApiKey.IsNull() {
		creds.apiKey = config.ApiKey.ValueString()
	}

	// Lookup env vars
	if creds.endpoint == "" {
		creds.endpoint = os.Getenv("STACKLET_ENDPOINT")
	}
	if creds.apiKey == "" {
		creds.apiKey = os.Getenv("STACKLET_API_KEY")
	}

	if homeDir, err := os.UserHomeDir(); err == nil {
		configFile := path.Join(homeDir, ".stacklet", "config.json")
		if content, err := os.ReadFile(configFile); err == nil {
			config := struct {
				Api string `json:"api"`
			}{}
			if err := json.Unmarshal(content, &config); err == nil {
				creds.endpoint = config.Api
			}
		}

		credsFile := path.Join(homeDir, ".stacklet", "credentials")
		if content, err := os.ReadFile(credsFile); err == nil {
			creds.apiKey = string(content)
		}
	}

	return &creds

}
