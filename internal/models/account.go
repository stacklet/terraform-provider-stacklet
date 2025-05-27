// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AccountResource is the model for account resources.
type AccountResource struct {
	ID                       types.String `tfsdk:"id"`
	Key                      types.String `tfsdk:"key"`
	Name                     types.String `tfsdk:"name"`
	ShortName                types.String `tfsdk:"short_name"`
	Description              types.String `tfsdk:"description"`
	CloudProvider            types.String `tfsdk:"cloud_provider"`
	Path                     types.String `tfsdk:"path"`
	Email                    types.String `tfsdk:"email"`
	SecurityContext          types.String `tfsdk:"security_context"`
	SecurityContextWO        types.String `tfsdk:"security_context_wo"`
	SecurityContextWOVersion types.String `tfsdk:"security_context_wo_version"`
	Variables                types.String `tfsdk:"variables"`
}

// AcountDataSourcemodel is the model for account data sources.
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
