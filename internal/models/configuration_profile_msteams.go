// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
)

// ConfigurationProfileMSTeamsDataSource is the model for Microsoft Teams configuration profile data sources.
type ConfigurationProfileMSTeamsDataSource struct {
	ID              types.String `tfsdk:"id"`
	Profile         types.String `tfsdk:"profile"`
	AccessConfig    types.Object `tfsdk:"access_config"`
	CustomerConfig  types.Object `tfsdk:"customer_config"`
	ChannelMappings types.List   `tfsdk:"channel_mapping"`
	EntityDetails   types.Object `tfsdk:"entity_details"`
}

func (c ConfigurationProfileMSTeamsDataSource) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"access_config": types.ObjectType{
			AttrTypes: MSTeamsAccessConfig{}.AttributeTypes(),
		},
		"customer_config": types.ObjectType{
			AttrTypes: MSTeamsCustomerConfig{}.AttributeTypes(),
		},
		"channel_mapping": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: MSTeamsChannelMapping{}.AttributeTypes(),
			},
		},
		"entity_details": types.ObjectType{
			AttrTypes: MSTeamsEntityDetails{}.AttributeTypes(),
		},
	}
}

func (m *ConfigurationProfileMSTeamsDataSource) Update(cp api.ConfigurationProfile) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(cp.ID)
	m.Profile = types.StringValue(cp.Profile)

	channelMappings, d := tftypes.ObjectList[MSTeamsChannelMapping](
		cp.Record.MSTeamsConfiguration.ChannelMappings,
		func(entry api.MSTeamsChannelMapping) (map[string]attr.Value, diag.Diagnostics) {
			return map[string]attr.Value{
				"name":       types.StringValue(entry.Name),
				"team_id":    types.StringValue(string(entry.TeamID)),
				"channel_id": types.StringValue(entry.ChannelID),
			}, nil
		},
	)
	diags.Append(d...)
	m.ChannelMappings = channelMappings

	accessConfig, d := m.buildAccessConfig(cp)
	diags.Append(d...)
	m.AccessConfig = accessConfig

	customerConfig, d := m.buildCustomerConfig(cp)
	diags.Append(d...)
	m.CustomerConfig = customerConfig

	entityDetails, d := m.buildEntityDetails(cp)
	diags.Append(d...)
	m.EntityDetails = entityDetails

	return diags
}

func (m ConfigurationProfileMSTeamsDataSource) buildAccessConfig(cp api.ConfigurationProfile) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	cfg := cp.Record.MSTeamsConfiguration.AccessConfig
	if cfg == nil {
		return types.ObjectNull(MSTeamsAccessConfig{}.AttributeTypes()), diags
	}

	botApplication, d := types.ObjectValue(
		MSTeamsBotApplication{}.AttributeTypes(),
		map[string]attr.Value{
			"download_url": types.StringValue(cfg.BotApplication.DownloadURL),
			"version":      types.StringValue(cfg.BotApplication.Version),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	var publishedApplication basetypes.ObjectValue
	if cfg.PublishedApplication == nil || (cfg.PublishedApplication.CatalogID == nil && cfg.PublishedApplication.Version == nil) {
		publishedApplication = types.ObjectNull(MSTeamsPublishedApplication{}.AttributeTypes())
	} else {
		var d diag.Diagnostics
		publishedApplication, d = types.ObjectValue(
			MSTeamsPublishedApplication{}.AttributeTypes(),
			map[string]attr.Value{
				"catalog_id": types.StringPointerValue(cfg.PublishedApplication.CatalogID),
				"version":    types.StringPointerValue(cfg.PublishedApplication.Version),
			},
		)
		diags.Append(d...)
		if diags.HasError() {
			return basetypes.ObjectValue{}, diags
		}
	}

	return types.ObjectValue(
		MSTeamsAccessConfig{}.AttributeTypes(),
		map[string]attr.Value{
			"client_id":             types.StringValue(cfg.ClientID),
			"roundtrip_digest":      types.StringValue(cfg.RoundtripDigest),
			"tenant_id":             types.StringValue(cfg.TenantID),
			"bot_application":       botApplication,
			"published_application": publishedApplication,
		},
	)
}

func (m ConfigurationProfileMSTeamsDataSource) buildCustomerConfig(cp api.ConfigurationProfile) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	cfg := cp.Record.MSTeamsConfiguration.CustomerConfig

	tags, d := cfg.Tags.TagsMap()
	diags.Append(d...)
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	var version types.String
	if cfg.TerraformModule.Version != nil {
		version = types.StringValue(*cfg.TerraformModule.Version)
	} else {
		version = types.StringNull()
	}

	terraformModule, d := types.ObjectValue(
		TerraformModule{}.AttributeTypes(),
		map[string]attr.Value{
			"repository_url": types.StringValue(cfg.TerraformModule.RepositoryURL),
			"source":         types.StringValue(cfg.TerraformModule.Source),
			"version":        version,
			"variables_json": types.StringValue(cfg.TerraformModule.VariablesJSON),
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	return types.ObjectValue(
		MSTeamsCustomerConfig{}.AttributeTypes(),
		map[string]attr.Value{
			"bot_endpoint":     types.StringValue(cfg.BotEndpoint),
			"oidc_client":      types.StringValue(cfg.OIDCClient),
			"oidc_issuer":      types.StringValue(cfg.OIDCIssuer),
			"prefix":           types.StringValue(cfg.Prefix),
			"roundtrip_digest": types.StringValue(cfg.RoundtripDigest),
			"tags":             tags,
			"terraform_module": terraformModule,
		},
	)
}

func (m ConfigurationProfileMSTeamsDataSource) buildEntityDetails(cp api.ConfigurationProfile) (basetypes.ObjectValue, diag.Diagnostics) {
	channels, diags := tftypes.ObjectList[MSTeamsChannelDetails](
		cp.Record.MSTeamsConfiguration.EntityDetails.Channels,
		func(entry api.MSTeamsChannelDetail) (map[string]attr.Value, diag.Diagnostics) {
			return map[string]attr.Value{
				"id":   types.StringValue(entry.ID),
				"name": types.StringValue(entry.Name),
			}, nil
		},
	)
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	teams, teamsDiags := tftypes.ObjectList[MSTeamsTeamDetails](
		cp.Record.MSTeamsConfiguration.EntityDetails.Teams,
		func(entry api.MSTeamsTeamDetail) (map[string]attr.Value, diag.Diagnostics) {
			return map[string]attr.Value{
				"id":   types.StringValue(entry.ID),
				"name": types.StringValue(entry.Name),
			}, nil
		},
	)
	diags.Append(teamsDiags...)
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	return types.ObjectValue(
		MSTeamsEntityDetails{}.AttributeTypes(),
		map[string]attr.Value{
			"channels": channels,
			"teams":    teams,
		},
	)
}

// ConfigurationProfileMSTeamsResource is the model for Microsoft Teams configuration profile resources.
type ConfigurationProfileMSTeamsResource struct {
	ConfigurationProfileMSTeamsDataSource

	AccessConfigInput   types.Object `tfsdk:"access_config_input"`
	CustomerConfigInput types.Object `tfsdk:"customer_config_input"`
}

func (r ConfigurationProfileMSTeamsResource) AttributeTypes() map[string]attr.Type {
	attrTypes := r.ConfigurationProfileMSTeamsDataSource.AttributeTypes()
	attrTypes["access_config_input"] = types.ObjectType{
		AttrTypes: MSTeamsAccessConfigInput{}.AttributeTypes(),
	}
	attrTypes["customer_config_input"] = types.ObjectType{
		AttrTypes: MSTeamsCustomerConfigInput{}.AttributeTypes(),
	}
	return attrTypes
}

func (m *ConfigurationProfileMSTeamsResource) Update(cp api.ConfigurationProfile) diag.Diagnostics {
	diags := m.ConfigurationProfileMSTeamsDataSource.Update(cp)

	// update state for input fields to match what's returned by the API.
	// Since the inputs are separated from the outputs, this is needed to match
	// the current state (including when the resource is imported).
	accessConfigInput, d := tftypes.FilteredObject[MSTeamsAccessConfigInput](
		m.AccessConfig,
		[]string{"client_id", "roundtrip_digest", "tenant_id"},
	)
	m.AccessConfigInput = accessConfigInput
	diags.Append(d...)

	customerConfigInput, d := tftypes.FilteredObject[MSTeamsCustomerConfigInput](
		m.CustomerConfig,
		[]string{"prefix", "tags"},
	)
	m.CustomerConfigInput = customerConfigInput
	diags.Append(d...)

	return diags
}

// MSTeamsAccessConfigInput is the model for Microsoft Teams access configuration input (user-provided fields).
type MSTeamsAccessConfigInput struct {
	ClientID        types.String `tfsdk:"client_id"`
	RoundtripDigest types.String `tfsdk:"roundtrip_digest"`
	TenantID        types.String `tfsdk:"tenant_id"`
}

func (a MSTeamsAccessConfigInput) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"client_id":        types.StringType,
		"roundtrip_digest": types.StringType,
		"tenant_id":        types.StringType,
	}
}

// MSTeamsCustomerConfigInput is the model for Microsoft Teams customer configuration input (user-provided fields).
type MSTeamsCustomerConfigInput struct {
	Prefix types.String `tfsdk:"prefix"`
	Tags   types.Map    `tfsdk:"tags"`
}

func (c MSTeamsCustomerConfigInput) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"prefix": types.StringType,
		"tags": types.MapType{
			ElemType: types.StringType,
		},
	}
}

// MSTeamsChannelMapping is the model for Microsoft Teams channel mappings.
type MSTeamsChannelMapping struct {
	Name      types.String `tfsdk:"name"`
	TeamID    types.String `tfsdk:"team_id"`
	ChannelID types.String `tfsdk:"channel_id"`
}

func (c MSTeamsChannelMapping) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":       types.StringType,
		"team_id":    types.StringType,
		"channel_id": types.StringType,
	}
}

// MSTeamsEntityDetails is the model for Microsoft Teams entity details.
type MSTeamsEntityDetails struct {
	Channels types.List `tfsdk:"channels"`
	Teams    types.List `tfsdk:"teams"`
}

func (e MSTeamsEntityDetails) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"channels": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: MSTeamsChannelDetails{}.AttributeTypes(),
			},
		},
		"teams": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: MSTeamsTeamDetails{}.AttributeTypes(),
			},
		},
	}
}

// MSTeamsChannelDetails is the model for Microsoft Teams channel details.
type MSTeamsChannelDetails struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d MSTeamsChannelDetails) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}
}

// MSTeamsTeamDetails is the model for Microsoft Teams team details.
type MSTeamsTeamDetails struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d MSTeamsTeamDetails) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}
}

// MSTeamsCustomerConfig is the model for Microsoft Teams customer configuration.
type MSTeamsCustomerConfig struct {
	BotEndpoint     types.String `tfsdk:"bot_endpoint"`
	OIDCClient      types.String `tfsdk:"oidc_client"`
	OIDCIssuer      types.String `tfsdk:"oidc_issuer"`
	Prefix          types.String `tfsdk:"prefix"`
	RoundtripDigest types.String `tfsdk:"roundtrip_digest"`
	Tags            types.Map    `tfsdk:"tags"`
	TerraformModule types.Object `tfsdk:"terraform_module"`
}

func (c MSTeamsCustomerConfig) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"bot_endpoint":     types.StringType,
		"oidc_client":      types.StringType,
		"oidc_issuer":      types.StringType,
		"prefix":           types.StringType,
		"roundtrip_digest": types.StringType,
		"tags": types.MapType{
			ElemType: types.StringType,
		},
		"terraform_module": types.ObjectType{
			AttrTypes: TerraformModule{}.AttributeTypes(),
		},
	}
}

// MSTeamsAccessConfig is the model for Microsoft Teams access configuration.
type MSTeamsAccessConfig struct {
	ClientID             types.String `tfsdk:"client_id"`
	RoundtripDigest      types.String `tfsdk:"roundtrip_digest"`
	TenantID             types.String `tfsdk:"tenant_id"`
	BotApplication       types.Object `tfsdk:"bot_application"`
	PublishedApplication types.Object `tfsdk:"published_application"`
}

func (a MSTeamsAccessConfig) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"client_id":        types.StringType,
		"roundtrip_digest": types.StringType,
		"tenant_id":        types.StringType,
		"bot_application": types.ObjectType{
			AttrTypes: MSTeamsBotApplication{}.AttributeTypes(),
		},
		"published_application": types.ObjectType{
			AttrTypes: MSTeamsPublishedApplication{}.AttributeTypes(),
		},
	}
}

// MSTeamsBotApplication is the model for Microsoft Teams bot application.
type MSTeamsBotApplication struct {
	DownloadURL types.String `tfsdk:"download_url"`
	Version     types.String `tfsdk:"version"`
}

func (b MSTeamsBotApplication) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"download_url": types.StringType,
		"version":      types.StringType,
	}
}

// MSTeamsPublishedApplication is the model for Microsoft Teams published application.
type MSTeamsPublishedApplication struct {
	CatalogID types.String `tfsdk:"catalog_id"`
	Version   types.String `tfsdk:"version"`
}

func (p MSTeamsPublishedApplication) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"catalog_id": types.StringType,
		"version":    types.StringType,
	}
}
