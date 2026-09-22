// Copyright Stacklet, Inc. 2025, 2026

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
)

// DatadogIntegrationID is the identifier for the Datadog integration.
//
// The integration applies to the whole deployment and has no identifier of its
// own in the API, so a fixed one stands in for it.
const DatadogIntegrationID = "datadog"

// DatadogIntegrationDataSource is the model for Datadog integration data sources.
type DatadogIntegrationDataSource struct {
	ID         types.String `tfsdk:"id"`
	Enabled    types.Bool   `tfsdk:"enabled"`
	Site       types.String `tfsdk:"site"`
	Configured types.Bool   `tfsdk:"configured"`
	Status     types.Object `tfsdk:"status"`
}

func (m *DatadogIntegrationDataSource) Update(integration *api.DatadogIntegration) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(DatadogIntegrationID)
	m.Enabled = types.BoolValue(integration.Enabled)
	m.Site = types.StringValue(integration.Site)
	m.Configured = types.BoolValue(integration.Configured)

	status, d := NewStatusInfo(integration.StatusInfo)
	m.Status = status
	errors.AddAttributeDiags(&diags, d, "status")

	return diags
}

// DatadogIntegrationResource is the model for Datadog integration resources.
//
// The credentials are write-only: the API never returns them, so nothing about
// them is kept in state beyond the versions that say when to send them again.
type DatadogIntegrationResource struct {
	DatadogIntegrationDataSource

	APIKeyWO        types.String `tfsdk:"api_key_wo"`
	APIKeyWOVersion types.String `tfsdk:"api_key_wo_version"`
	AppKeyWO        types.String `tfsdk:"app_key_wo"`
	AppKeyWOVersion types.String `tfsdk:"app_key_wo_version"`
}
