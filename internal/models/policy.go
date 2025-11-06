// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
)

// PolicyDataSource is the model for policy data sources.
type PolicyDataSource struct {
	ID              types.String `tfsdk:"id"`
	UUID            types.String `tfsdk:"uuid"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	CloudProvider   types.String `tfsdk:"cloud_provider"`
	Version         types.Int32  `tfsdk:"version"`
	Category        types.List   `tfsdk:"category"`
	Mode            types.String `tfsdk:"mode"`
	ResourceType    types.String `tfsdk:"resource_type"`
	Path            types.String `tfsdk:"path"`
	SourceJSON      types.String `tfsdk:"source_json"`
	SourceYAML      types.String `tfsdk:"source_yaml"`
	System          types.Bool   `tfsdk:"system"`
	UnqualifiedName types.String `tfsdk:"unqualified_name"`
}

func (m *PolicyDataSource) Update(ctx context.Context, policy *api.Policy) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(policy.ID)
	m.UUID = types.StringValue(policy.UUID)
	m.Name = types.StringValue(policy.Name)
	m.Description = types.StringPointerValue(policy.Description)
	m.CloudProvider = types.StringValue(policy.Provider)
	m.Version = types.Int32Value(int32(policy.Version))

	m.Mode = types.StringValue(policy.Mode)
	m.ResourceType = types.StringValue(policy.ResourceType)
	m.Path = types.StringValue(policy.Path)
	m.SourceJSON = types.StringValue(policy.Source)
	m.SourceYAML = types.StringValue(policy.SourceYAML)
	m.System = types.BoolValue(policy.System)
	m.UnqualifiedName = types.StringValue(policy.UnqualifiedName)

	category, d := types.ListValueFrom(ctx, types.StringType, policy.Category)
	m.Category = category
	errors.AddAttributeDiags(&diags, d, "category")

	return diags
}
