// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

// AccountDiscoveryAzureResource is the model for Azure account discovery resources.
type AccountDiscoveryAzureResource struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Suspended    types.Bool   `tfsdk:"suspended"`
	ClientID     types.String `tfsdk:"client_id"`
	TenantID     types.String `tfsdk:"tenant_id"`
	ClientSecret types.String `tfsdk:"client_secret_wo"`
}

func (m *AccountDiscoveryAzureResource) Update(accountDiscovery *api.AccountDiscovery) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(accountDiscovery.ID)
	m.Name = types.StringValue(accountDiscovery.Name)
	m.Description = types.StringPointerValue(accountDiscovery.Description)
	m.Suspended = types.BoolValue(accountDiscovery.Schedule.Suspended)
	m.ClientID = types.StringValue(accountDiscovery.Config.AzureConfig.ClientID)
	m.TenantID = types.StringValue(accountDiscovery.Config.AzureConfig.TenantID)

	return diags
}
