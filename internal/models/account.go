// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/typehelpers"
)

// AcountDataSource is the model for account data sources.
type AccountDataSource struct {
	ID              types.String `tfsdk:"id"`
	Key             types.String `tfsdk:"key"`
	Name            types.String `tfsdk:"name"`
	ShortName       types.String `tfsdk:"short_name"`
	Description     types.String `tfsdk:"description"`
	CloudProvider   types.String `tfsdk:"cloud_provider"`
	Path            types.String `tfsdk:"path"`
	Email           types.String `tfsdk:"email"`
	SecurityContext types.String `tfsdk:"security_context"`
	Variables       types.String `tfsdk:"variables"`
}

func (m *AccountDataSource) Update(account *api.Account) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(account.ID)
	m.Key = types.StringValue(account.Key)
	m.Name = types.StringValue(account.Name)
	m.ShortName = types.StringPointerValue(account.ShortName)
	m.Description = types.StringPointerValue(account.Description)
	m.CloudProvider = types.StringValue(string(account.Provider))
	m.Path = types.StringPointerValue(account.Path)
	m.Email = types.StringPointerValue(account.Email)
	m.SecurityContext = types.StringPointerValue(account.SecurityContext)

	variablesString, d := typehelpers.JSONString(account.Variables)
	errors.AddAttributeDiags(&diags, d, "variables")
	m.Variables = variablesString

	return diags
}

// AccountResource is the model for account resources.
type AccountResource struct {
	AccountDataSource

	SecurityContextWO        types.String `tfsdk:"security_context_wo"`
	SecurityContextWOVersion types.String `tfsdk:"security_context_wo_version"`
}
