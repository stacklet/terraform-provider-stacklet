// Copyright Stacklet, Inc. 2025, 2026

package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api/iac"
	"github.com/stacklet/terraform-provider-stacklet/internal/typehelpers"
)

// IaCPoliciesDataSource is the model for the IaC policies data source.
type IaCPoliciesDataSource struct {
	Source     types.String `tfsdk:"source"`
	Collection types.String `tfsdk:"collection"`
	Policies   types.List   `tfsdk:"policies"`
}

func (m *IaCPoliciesDataSource) Update(ctx context.Context, policies []iac.Policy) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Policies, diags = typehelpers.ObjectList[IaCPolicy](
		policies,
		func(policy iac.Policy) (map[string]attr.Value, diag.Diagnostics) {
			return map[string]attr.Value{
				"name":      types.StringValue(policy.Name),
				"source":    types.StringValue(policy.SourceName),
				"resources": typehelpers.StringsList(policy.Resources),
				"severity":  types.StringValue(policy.Severity),
			}, nil
		},
	)

	return diags
}

// IaCPolicy is an IaC policy in the policies data source.
type IaCPolicy struct {
	Name      types.String `tfsdk:"name"`
	Source    types.String `tfsdk:"source"`
	Resources types.List   `tfsdk:"resources"`
	Severity  types.String `tfsdk:"severity"`
}

func (p IaCPolicy) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":      types.StringType,
		"source":    types.StringType,
		"resources": types.ListType{ElemType: types.StringType},
		"severity":  types.StringType,
	}
}
