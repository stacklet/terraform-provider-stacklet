// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
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
