// Copyright Stacklet, Inc. 2025, 2026

package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api/iac"
)

// IaCPolicyCollectionBase holds attributes common to the policy collection resource and data source.
type IaCPolicyCollectionBase struct {
	Name      types.String `tfsdk:"name"`
	Policies  types.Set    `tfsdk:"policies"`
	IsDefault types.Bool   `tfsdk:"is_default"`
}

func (m *IaCPolicyCollectionBase) Update(ctx context.Context, collection *iac.PolicyCollection) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Name = types.StringValue(collection.Name)
	m.IsDefault = types.BoolValue(collection.IsDefault)
	policies := collection.Policies
	if policies == nil {
		policies = []string{}
	}
	m.Policies, diags = types.SetValueFrom(ctx, types.StringType, policies)

	return diags
}

type IaCPolicyCollectionDataSource struct {
	IaCPolicyCollectionBase

	IsMirror types.Bool `tfsdk:"is_mirror"`
}

func (m *IaCPolicyCollectionDataSource) Update(ctx context.Context, collection *iac.PolicyCollection) diag.Diagnostics {
	diags := m.IaCPolicyCollectionBase.Update(ctx, collection)
	m.IsMirror = types.BoolValue(collection.IsMirror)
	return diags
}

type IaCPolicyCollectionResource struct {
	IaCPolicyCollectionBase
}
