// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/typehelpers"
)

// AccountGroupResource is the model for account group resources.
type AccountGroupResource struct {
	ID            types.String `tfsdk:"id"`
	UUID          types.String `tfsdk:"uuid"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	CloudProvider types.String `tfsdk:"cloud_provider"`
	Regions       types.List   `tfsdk:"regions"`
}

func (m *AccountGroupResource) Update(accountGroup *api.AccountGroup) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(accountGroup.ID)
	m.UUID = types.StringValue(accountGroup.UUID)
	m.Name = types.StringValue(accountGroup.Name)
	m.Description = types.StringPointerValue(accountGroup.Description)
	m.CloudProvider = types.StringValue(accountGroup.Provider)
	m.Regions = typehelpers.StringsList(accountGroup.Regions)

	return diags
}

// AccountGroupDataSource is the model for account group data sources.
type AccountGroupDataSource AccountGroupResource
