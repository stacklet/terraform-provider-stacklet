// Copyright Stacklet, Inc. 2025, 2026

package provider

import (
	"context"
	"encoding/json"
	"os"
	"path"

	"github.com/caarlos0/env/v11"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/datasources"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
	"github.com/stacklet/terraform-provider-stacklet/internal/resources"
)

var _ provider.Provider = &stackletProvider{}

// providerModel holds the terraform configuration for the provider.
type providerModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	APIKey   types.String `tfsdk:"api_key"`
}

// providerEnv holds environment variables supported by the provider.
type providerEnv struct {
	Endpoint           string `env:"STACKLET_ENDPOINT"`
	APIKey             string `env:"STACKLET_API_KEY"`
	PageSize           int    `env:"STACKLET_PAGE_SIZE" envDefault:"100"`
	UnreleasedFeatures bool   `env:"STACKLET_UNRELEASED_FEATURES"`
}

type stackletProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and run locally, and "test" when running acceptance
	// testing.
	version string
}

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
		Description: `
This provider interacts with Stacklet's cloud governance platform.

It allows managing resources like accounts, account groups, policy collections, bindings and so on.
`,
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Description: `
The API key for Stacklet authentication.

May also be provided via STACKLET_API_KEY environment variable, or from the stacklet-admin CLI configuration.
`,
				Optional:  true,
				Sensitive: true,
			},
			"endpoint": schema.StringAttribute{
				Description: `
The endpoint URL of the Stacklet GraphQL API.

 May also be provided via STACKLET_ENDPOINT environment variable, or from the stacklet-admin CLI configuration.
`,
				Optional: true,
			},
		},
	}
}

// Configure prepares a Stacklet API client for data sources and resources.
func (p *stackletProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config providerModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	env, err := envConfig()
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	creds, diags := getCredentials(config, env)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Make provider data accessible to the Configure method of resources and data sources
	providerData := providerdata.New(
		ctx,
		api.ClientConfig{
			Endpoint: creds.Endpoint,
			APIKey:   creds.APIKey,
			Version:  p.version,
			PageSize: env.PageSize,
		},
	)
	resp.ResourceData = providerData
	resp.DataSourceData = providerData
}

// DataSources defines the data sources implemented in the provider.
func (p *stackletProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	conf, _ := envConfig()
	return datasources.DataSources.List(conf.UnreleasedFeatures)
}

// Resources defines the resources implemented in the provider.
func (p *stackletProvider) Resources(_ context.Context) []func() resource.Resource {
	conf, _ := envConfig()
	return resources.Resources.List(conf.UnreleasedFeatures)
}

func envConfig() (providerEnv, error) {
	return env.ParseAs[providerEnv]()
}

type credentials struct {
	Endpoint string
	APIKey   string
}

func getCredentials(config providerModel, env providerEnv) (credentials, diag.Diagnostics) {
	var creds credentials
	var diags diag.Diagnostics

	if config.Endpoint.IsUnknown() {
		diags.AddAttributeError(
			tfpath.Root("endpoint"),
			"Unknown Stacklet API Endpoint",
			"The provider cannot create the Stacklet API client as there is an unknown configuration value for the Stacklet API endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the STACKLET_ENDPOINT environment variable.",
		)
	}
	if config.APIKey.IsUnknown() {
		diags.AddAttributeError(
			tfpath.Root("api_key"),
			"Unknown Stacklet API key",
			"The provider cannot create the Stacklet API client as there is an unknown configuration value for the Stacklet API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the STACKLET_API_KEY environment variable.",
		)
	}
	if diags.HasError() {
		return creds, diags
	}

	// lookup provider configuration (might return empty strings)
	creds.Endpoint = config.Endpoint.ValueString()
	creds.APIKey = config.APIKey.ValueString()

	// lookup environment variables
	if creds.Endpoint == "" {
		creds.Endpoint = env.Endpoint
	}
	if creds.APIKey == "" {
		creds.APIKey = env.APIKey
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

	if creds.Endpoint == "" {
		diags.AddAttributeError(
			tfpath.Root("endpoint"),
			"Missing Stacklet API Endpoint",
			"The provider cannot create the Stacklet API client as there is a missing or empty value for the Stacklet API endpoint. "+
				"Set the endpoint value in the configuration, in the STACKLET_ENDPOINT environment variable, or login via the stacklet-admin CLI first.",
		)
	}
	if creds.APIKey == "" {
		diags.AddAttributeError(
			tfpath.Root("api_key"),
			"Missing Stacklet API key",
			"The provider cannot create the Stacklet API client as there is a missing or empty value for the Stacklet API key. "+
				"Set the api_key value in the configuration, in the STACKLET_API_KEY environment variable, or login via the stacklet-admin CLI first.",
		)
	}
	return creds, diags
}
