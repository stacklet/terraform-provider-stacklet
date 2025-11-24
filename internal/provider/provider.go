// Copyright (c) 2025 - Stacklet, Inc.

package provider

import (
	"context"
	"encoding/json"
	"os"
	"path"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	tfpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/datasources"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
	"github.com/stacklet/terraform-provider-stacklet/internal/resources"
)

var (
	_ provider.Provider = &stackletProvider{}
)

type stackletProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	APIKey   types.String `tfsdk:"api_key"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &stackletProvider{
			version: version,
		}
	}
}

type stackletProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and run locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *stackletProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "stacklet"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *stackletProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
This provider interacts with Stacklet's cloud governance platform.

It allows managing resources like accounts, account groups, policy collections, bindings and so on.
`,
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Description: `
The endpoint URL of the Stacklet GraphQL API.

 May also be provided via STACKLET_ENDPOINT environment variable, or from the stacklet-admin CLI configuration.
`,
				Optional: true,
			},
			"api_key": schema.StringAttribute{
				Description: `
The API key for Stacklet authentication.

May also be provided via STACKLET_API_KEY environment variable, or from the stacklet-admin CLI configuration.
`,
				Optional:  true,
				Sensitive: true,
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
	if config.APIKey.IsUnknown() {
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
	if creds.Endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			tfpath.Root("endpoint"),
			"Missing Stacklet API Endpoint",
			"The provider cannot create the Stacklet API client as there is a missing or empty value for the Stacklet API endpoint. "+
				"Set the endpoint value in the configuration, in the STACKLET_ENDPOINT environment variable, or login via the stacklet-admin CLI first.",
		)
	}
	if creds.APIKey == "" {
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

	// Make provider data accessible to the Configure method of resources and data sources
	providerData := providerdata.New(api.NewClient(ctx, creds.Endpoint, creds.APIKey))
	resp.ResourceData = providerData
	resp.DataSourceData = providerData

	tflog.Info(ctx, "Configured Stacklet client")
}

// DataSources defines the data sources implemented in the provider.
func (p *stackletProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return datasources.DataSources()
}

// Resources defines the resources implemented in the provider.
func (p *stackletProvider) Resources(_ context.Context) []func() resource.Resource {
	return resources.Resources()
}

type credentials struct {
	Endpoint string
	APIKey   string
}

func getCredentials(config stackletProviderModel) *credentials {
	creds := credentials{}

	// lookup provider configuration (might return empty strings)
	creds.Endpoint = config.Endpoint.ValueString()
	creds.APIKey = config.APIKey.ValueString()

	// lookup environment variables
	if creds.Endpoint == "" {
		creds.Endpoint = os.Getenv("STACKLET_ENDPOINT")
	}
	if creds.APIKey == "" {
		creds.APIKey = os.Getenv("STACKLET_API_KEY")
	}

	// lookup stacklet-admin configuration
	if homeDir, err := os.UserHomeDir(); err == nil {

		if creds.Endpoint == "" {
			configFile := path.Join(homeDir, ".stacklet", "config.json")
			if content, err := os.ReadFile(configFile); err == nil {
				config := struct {
					Api string `json:"api"`
				}{}
				if err := json.Unmarshal(content, &config); err == nil {
					creds.Endpoint = config.Api
				}
			}
		}

		if creds.APIKey == "" {
			credsFile := path.Join(homeDir, ".stacklet", "credentials")
			if content, err := os.ReadFile(credsFile); err == nil {
				creds.APIKey = string(content)
			}
		}
	}

	return &creds
}
