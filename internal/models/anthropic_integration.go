// Copyright Stacklet, Inc. 2025, 2026

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/typehelpers"
)

// AnthropicIntegrationDataSource is the model for Anthropic integration data sources.
type AnthropicIntegrationDataSource struct {
	ID         types.String `tfsdk:"id"`
	OrgID      types.String `tfsdk:"org_id"`
	Configured types.Bool   `tfsdk:"configured"`
}

func (m *AnthropicIntegrationDataSource) Update(integration *api.AnthropicIntegration) diag.Diagnostics {
	m.ID = typehelpers.GraphQLIDValue(integration.ID)
	m.OrgID = types.StringValue(integration.OrgID)
	m.Configured = types.BoolValue(integration.Configured)
	return nil
}

// AnthropicIntegrationResource is the model for Anthropic integration resources.
//
// The admin key is write-only: the API never returns it, so nothing about it
// is kept in state beyond the version that says when to send it again.
type AnthropicIntegrationResource struct {
	AnthropicIntegrationDataSource

	AdminKeyWO        types.String `tfsdk:"admin_key_wo"`
	AdminKeyWOVersion types.String `tfsdk:"admin_key_wo_version"`
}
