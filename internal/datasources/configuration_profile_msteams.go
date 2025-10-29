// Copyright (c) 2025 - Stacklet, Inc.

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/modelupdate"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
)

var (
	_ datasource.DataSource = &configurationProfileMSTeamsDataSource{}
)

func NewConfigurationProfileMSTeamsDataSource() datasource.DataSource {
	return &configurationProfileMSTeamsDataSource{}
}

type configurationProfileMSTeamsDataSource struct {
	api *api.API
}

func (d *configurationProfileMSTeamsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configuration_profile_msteams"
}

func (d *configurationProfileMSTeamsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve information about the Microsoft Teams configuration profile.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the configuration profile.",
				Computed:    true,
			},
			"profile": schema.StringAttribute{
				Description: "The profile name.",
				Computed:    true,
			},
			"access_config": schema.SingleNestedAttribute{
				Description: "Access configuration for Microsoft Teams.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"client_id": schema.StringAttribute{
						Description: "The client ID.",
						Computed:    true,
					},
					"roundtrip_digest": schema.StringAttribute{
						Description: "The roundtrip digest.",
						Computed:    true,
					},
					"tenant_id": schema.StringAttribute{
						Description: "The tenant ID.",
						Computed:    true,
					},
					"bot_application": schema.SingleNestedAttribute{
						Description: "Bot application configuration.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"download_url": schema.StringAttribute{
								Description: "The bot application download URL.",
								Computed:    true,
							},
							"version": schema.StringAttribute{
								Description: "The bot application version.",
								Computed:    true,
							},
						},
					},
					"published_application": schema.SingleNestedAttribute{
						Description: "Published application configuration.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"catalog_id": schema.StringAttribute{
								Description: "The catalog ID.",
								Computed:    true,
							},
							"version": schema.StringAttribute{
								Description: "The published application version.",
								Computed:    true,
							},
						},
					},
				},
			},
			"entity_details": schema.SingleNestedAttribute{
				Description: "Entity details for Microsoft Teams.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"channels": schema.ListNestedAttribute{
						Description: "Channel details.",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description: "The channel ID.",
									Computed:    true,
								},
								"name": schema.StringAttribute{
									Description: "The channel name.",
									Computed:    true,
								},
							},
						},
					},
					"teams": schema.ListNestedAttribute{
						Description: "Team details.",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description: "The team ID.",
									Computed:    true,
								},
								"name": schema.StringAttribute{
									Description: "The team name.",
									Computed:    true,
								},
							},
						},
					},
				},
			},
			"customer_config": schema.SingleNestedAttribute{
				Description: "Customer configuration for Microsoft Teams.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"bot_endpoint": schema.StringAttribute{
						Description: "The bot endpoint URL.",
						Computed:    true,
					},
					"oidc_client": schema.StringAttribute{
						Description: "The OIDC client ID.",
						Computed:    true,
					},
					"oidc_issuer": schema.StringAttribute{
						Description: "The OIDC issuer URL.",
						Computed:    true,
					},
					"prefix": schema.StringAttribute{
						Description: "The prefix for resources.",
						Computed:    true,
					},
					"roundtrip_digest": schema.StringAttribute{
						Description: "The roundtrip digest.",
						Computed:    true,
					},
					"tags": schema.MapAttribute{
						Description: "Tags for the configuration as key-value pairs.",
						ElementType: types.StringType,
						Computed:    true,
					},
					"terraform_module": schema.SingleNestedAttribute{
						Description: "Terraform module configuration.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"repository_url": schema.StringAttribute{
								Description: "The repository URL.",
								Computed:    true,
							},
							"source": schema.StringAttribute{
								Description: "The module source.",
								Computed:    true,
							},
							"version": schema.StringAttribute{
								Description: "The module version.",
								Computed:    true,
							},
							"variables_json": schema.StringAttribute{
								Description: "The module variables as JSON.",
								Computed:    true,
							},
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"channel_mapping": schema.ListNestedBlock{
				Description: "Channel mapping configuration.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The mapping name.",
							Computed:    true,
						},
						"team_id": schema.StringAttribute{
							Description: "The team ID.",
							Computed:    true,
						},
						"channel_id": schema.StringAttribute{
							Description: "The channel ID.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *configurationProfileMSTeamsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if pd, err := providerdata.GetDataSourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		d.api = pd.API
	}
}

func (d *configurationProfileMSTeamsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.ConfigurationProfileMSTeamsDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, err := d.api.ConfigurationProfile.ReadMSTeams(ctx)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	data.ID = types.StringValue(config.ID)
	data.Profile = types.StringValue(config.Profile)

	updater := modelupdate.NewConfigurationProfileUpdater(*config)

	accessConfig, diags := updater.MSTeamsAccessConfig()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.AccessConfig = accessConfig

	customerConfig, diags := updater.MSTeamsCustomerConfig()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.CustomerConfig = customerConfig

	channelMappings, diags := updater.MSTeamsChannelMappings()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.ChannelMappings = channelMappings

	entityDetails, diags := updater.MSTeamsEntityDetails()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.EntityDetails = entityDetails

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
