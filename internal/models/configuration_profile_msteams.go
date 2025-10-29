// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
