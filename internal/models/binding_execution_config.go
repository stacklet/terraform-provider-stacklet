// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BindingExecutionConfigResource is the model for a binding execution config resource.
type BindingExecutionConfigResource struct {
	ID          types.String `tfsdk:"id"`
	BindingUUID types.String `tfsdk:"binding_uuid"`
	DryRun      types.Bool   `tfsdk:"dry_run"`
	Variables   types.String `tfsdk:"variables"`
}

// BindingExecutionConfigDataSource is the model for a binding execution config data source.
type BindingExecutionConfigDataSource BindingExecutionConfigResource
