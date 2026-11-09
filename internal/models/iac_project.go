// Copyright Stacklet, Inc. 2025, 2026

package models

import (
	"context"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api/iac"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
)

// IaCProjectDataSource is the model for IaC project data sources.
type IaCProjectDataSource struct {
	Name        types.String `tfsdk:"name"`
	ClientID    types.String `tfsdk:"client_id"`
	Collections types.Set    `tfsdk:"collections"`
	ReadGroups  types.Set    `tfsdk:"read_groups"`
}

func (m *IaCProjectDataSource) Update(ctx context.Context, project *iac.Project) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Name = types.StringValue(project.Name)
	m.ClientID = types.StringPointerValue(project.ClientID)

	collections, d := types.SetValueFrom(ctx, types.StringType, nonNil(project.Collections))
	m.Collections = collections
	errors.AddAttributeDiags(&diags, d, "collections")

	var readGroups []string
	if project.Groups != nil {
		readGroups = project.Groups.Read
	}
	groups, d := types.SetValueFrom(ctx, types.StringType, nonNil(readGroups))
	m.ReadGroups = groups
	errors.AddAttributeDiags(&diags, d, "read_groups")

	return diags
}

// IaCProjectResource is the model for IaC project resources.
type IaCProjectResource struct {
	IaCProjectDataSource
}

// ToAPIInput returns the API input for creating or updating the project.
func (m *IaCProjectResource) ToAPIInput(ctx context.Context) (iac.ProjectInput, diag.Diagnostics) {
	var diags diag.Diagnostics

	collections, d := sortedStrings(ctx, m.Collections)
	errors.AddAttributeDiags(&diags, d, "collections")
	readGroups, d := sortedStrings(ctx, m.ReadGroups)
	errors.AddAttributeDiags(&diags, d, "read_groups")

	return iac.ProjectInput{
		Collections: collections,
		Groups:      iac.ProjectGroups{Read: readGroups},
	}, diags
}

// nonNil returns an empty slice if s is nil, so that it's converted to an
// empty set rather than a null one.
func nonNil(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}

// sortedStrings returns the sorted elements of a set of strings, always
// non-nil so that they're sent to the API as an empty list.
func sortedStrings(ctx context.Context, s types.Set) ([]string, diag.Diagnostics) {
	values := []string{}
	if s.IsUnknown() {
		var diags diag.Diagnostics
		diags.AddError("Unknown value", "The set value is not known yet.")
		return nil, diags
	}
	if s.IsNull() {
		return values, nil
	}
	diags := s.ElementsAs(ctx, &values, false)
	slices.Sort(values)
	return values, diags
}
