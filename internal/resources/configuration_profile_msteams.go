// Copyright (c) 2025 - Stacklet, Inc.

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/modelupdate"
	"github.com/stacklet/terraform-provider-stacklet/internal/planmodifiers"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
	"github.com/stacklet/terraform-provider-stacklet/internal/schemavalidate"
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
)

var (
	_ resource.Resource                = &configurationProfileMSTeamsResource{}
	_ resource.ResourceWithConfigure   = &configurationProfileMSTeamsResource{}
	_ resource.ResourceWithImportState = &configurationProfileMSTeamsResource{}
)

func NewConfigurationProfileMSTeamsResource() resource.Resource {
	return &configurationProfileMSTeamsResource{}
}

type configurationProfileMSTeamsResource struct {
	api *api.API
}

func (r *configurationProfileMSTeamsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configuration_profile_msteams"
}

func (r *configurationProfileMSTeamsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Manage the Microsoft Teams configuration profile.

The profile is global, adding multiple resources of this kind will cause them to override each other.
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the configuration profile.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"profile": schema.StringAttribute{
				Description: "The profile name.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
			"access_config_input": schema.SingleNestedAttribute{
				Description: "Access configuration input for Microsoft Teams.",
				Optional:    true,
				PlanModifiers: []planmodifier.Object{
					planmodifiers.RequiresReplaceIfUnset(),
				},
				Attributes: map[string]schema.Attribute{
					"client_id": schema.StringAttribute{
						Description: "The client ID.",
						Required:    true,
					},
					"roundtrip_digest": schema.StringAttribute{
						Description: "The roundtrip digest.",
						Required:    true,
					},
					"tenant_id": schema.StringAttribute{
						Description: "The tenant ID.",
						Required:    true,
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
			"customer_config_input": schema.SingleNestedAttribute{
				Description: "Customer configuration input for Microsoft Teams.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					planmodifiers.RequiresReplaceIfUnset(),
					planmodifiers.DefaultObject(
						types.ObjectValueMust(
							map[string]attr.Type{
								"prefix": types.StringType,
								"tags": types.MapType{
									ElemType: types.StringType,
								},
							},
							map[string]attr.Value{
								"prefix": types.StringValue("stacklet"),
								"tags":   types.MapValueMust(types.StringType, map[string]attr.Value{}),
							},
						),
					),
				},
				Attributes: map[string]schema.Attribute{
					"prefix": schema.StringAttribute{
						Description: "The prefix for resources.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("stacklet"),
					},
					"tags": schema.MapAttribute{
						Description: "Tags for the configuration as key-value pairs.",
						ElementType: types.StringType,
						Optional:    true,
						Computed:    true,
						Default:     tftypes.EmptyMapDefault(types.StringType),
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
		},
		Blocks: map[string]schema.Block{
			"channel_mapping": schema.ListNestedBlock{
				Description: "Channel mapping configuration.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The mapping name.",
							Required:    true,
						},
						"team_id": schema.StringAttribute{
							Description: "The team ID.",
							Required:    true,
							Validators: []validator.String{
								schemavalidate.UUID(),
							},
						},
						"channel_id": schema.StringAttribute{
							Description: "The channel ID.",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func (r *configurationProfileMSTeamsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if pd, err := providerdata.GetResourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		r.api = pd.API
	}
}

func (r *configurationProfileMSTeamsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.ConfigurationProfileMSTeamsResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input, diags := r.buildInput(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, err := r.api.ConfigurationProfile.UpsertMSTeams(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(r.updateStateFromAPI(&data, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *configurationProfileMSTeamsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.ConfigurationProfileMSTeamsResource
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, err := r.api.ConfigurationProfile.ReadMSTeams(ctx)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(r.updateStateFromAPI(&data, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *configurationProfileMSTeamsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data models.ConfigurationProfileMSTeamsResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input, diags := r.buildInput(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, err := r.api.ConfigurationProfile.UpsertMSTeams(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(r.updateStateFromAPI(&data, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *configurationProfileMSTeamsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.ConfigurationProfileMSTeamsResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.api.ConfigurationProfile.DeleteMSTeams(ctx); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *configurationProfileMSTeamsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("profile"), string(api.ConfigurationProfileMSTeams))...)
}

func (r *configurationProfileMSTeamsResource) buildInput(ctx context.Context, data models.ConfigurationProfileMSTeamsResource) (api.MSTeamsConfigurationInput, diag.Diagnostics) {
	var diags diag.Diagnostics
	input := api.MSTeamsConfigurationInput{}

	if !data.AccessConfigInput.IsNull() && !data.AccessConfigInput.IsUnknown() {
		var accessConfigInput models.MSTeamsAccessConfigInput
		diags.Append(data.AccessConfigInput.As(ctx, &accessConfigInput, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return input, diags
		}

		input.AccessConfig = &api.MSTeamsAccessConfigInput{
			ClientID:        accessConfigInput.ClientID.ValueString(),
			RoundtripDigest: accessConfigInput.RoundtripDigest.ValueString(),
			TenantID:        accessConfigInput.TenantID.ValueString(),
		}
	}

	if !data.CustomerConfigInput.IsNull() && !data.CustomerConfigInput.IsUnknown() {
		var customerConfigInput models.MSTeamsCustomerConfigInput
		diags.Append(data.CustomerConfigInput.As(ctx, &customerConfigInput, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return input, diags
		}

		input.CustomerConfig = &api.MSTeamsCustomerConfigInput{
			Prefix: customerConfigInput.Prefix.ValueString(),
			Tags:   api.NewTagsList(customerConfigInput.Tags),
		}
	}

	if !data.ChannelMappings.IsNull() && !data.ChannelMappings.IsUnknown() {
		var channelMappings []models.MSTeamsChannelMapping
		diags.Append(data.ChannelMappings.ElementsAs(ctx, &channelMappings, false)...)
		if diags.HasError() {
			return input, diags
		}

		apiMappings := make([]api.MSTeamsChannelMapping, len(channelMappings))
		for i, mapping := range channelMappings {
			apiMappings[i] = api.MSTeamsChannelMapping{
				Name:      mapping.Name.ValueString(),
				TeamID:    api.UUID(mapping.TeamID.ValueString()),
				ChannelID: mapping.ChannelID.ValueString(),
			}
		}
		input.ChannelMappings = &apiMappings
	}

	return input, diags
}

func (r *configurationProfileMSTeamsResource) updateStateFromAPI(data *models.ConfigurationProfileMSTeamsResource, config *api.ConfigurationProfile) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(config.ID)
	data.Profile = types.StringValue(config.Profile)

	updater := modelupdate.NewConfigurationProfileUpdater(*config)

	accessConfig, d := updater.MSTeamsAccessConfig()
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}
	data.AccessConfig = accessConfig

	customerConfig, d := updater.MSTeamsCustomerConfig()
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}
	data.CustomerConfig = customerConfig

	channelMappings, d := updater.MSTeamsChannelMappings()
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}
	data.ChannelMappings = channelMappings

	entityDetails, d := updater.MSTeamsEntityDetails()
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}
	data.EntityDetails = entityDetails

	// update state for input fields to match what's returned by the API. Since
	// the inputs are separated from the outputs, this is needed to see diffs.
	accessConfigInput, d := r.fillInputFromOutput(
		accessConfig,
		models.MSTeamsAccessConfigInput{},
		[]string{"client_id", "roundtrip_digest", "tenant_id"},
	)
	data.AccessConfigInput = accessConfigInput
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	customerConfigInput, d := r.fillInputFromOutput(
		customerConfig,
		models.MSTeamsCustomerConfigInput{},
		[]string{"prefix", "tags"},
	)
	data.CustomerConfigInput = customerConfigInput
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	return diags
}

func (r *configurationProfileMSTeamsResource) fillInputFromOutput(o basetypes.ObjectValue, t tftypes.WithAttributes, fields []string) (basetypes.ObjectValue, diag.Diagnostics) {
	if o.IsNull() || o.IsUnknown() {
		return tftypes.NullObject(t), nil
	}

	attrs := o.Attributes()
	values := make(map[string]attr.Value)
	// fill relevant fields
	for _, field := range fields {
		values[field] = attrs[field]
	}
	return types.ObjectValue(t.AttributeTypes(), values)
}
