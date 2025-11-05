// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/typehelpers"
)

// JiraProject is the model for a Jira project.
type JiraProject struct {
	ClosedStatus types.String `tfsdk:"closed_status"`
	IssueType    types.String `tfsdk:"issue_type"`
	Name         types.String `tfsdk:"name"`
	Project      types.String `tfsdk:"project"`
}

func (p JiraProject) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"closed_status": types.StringType,
		"issue_type":    types.StringType,
		"name":          types.StringType,
		"project":       types.StringType,
	}
}

// ConfigurationProfileJiraDataSource is the model for Jira configuration profile data sources.
type ConfigurationProfileJiraDataSource struct {
	ID       types.String `tfsdk:"id"`
	Profile  types.String `tfsdk:"profile"`
	URL      types.String `tfsdk:"url"`
	Projects types.List   `tfsdk:"project"`
	User     types.String `tfsdk:"user"`
	APIKey   types.String `tfsdk:"api_key"`
}

func (m *ConfigurationProfileJiraDataSource) Update(cp api.ConfigurationProfile) diag.Diagnostics {
	var diags diag.Diagnostics

	jiraConfig := cp.Record.JiraConfiguration

	m.ID = types.StringValue(cp.ID)
	m.Profile = types.StringValue(cp.Profile)
	m.URL = types.StringPointerValue(jiraConfig.URL)
	m.User = types.StringValue(jiraConfig.User)
	m.APIKey = types.StringValue(jiraConfig.APIKey)

	projects, d := typehelpers.ObjectList[JiraProject](
		jiraConfig.Projects,
		func(entry api.JiraProject) (map[string]attr.Value, diag.Diagnostics) {
			return map[string]attr.Value{
				"closed_status": types.StringValue(entry.ClosedStatus),
				"issue_type":    types.StringValue(entry.IssueType),
				"name":          types.StringValue(entry.Name),
				"project":       types.StringValue(entry.Project),
			}, nil
		},
	)
	m.Projects = projects
	diags.Append(d...)

	return diags
}

// ConfigurationProfileJiraResource is the model for Jira configuration profile resources.
type ConfigurationProfileJiraResource struct {
	ConfigurationProfileJiraDataSource

	APIKeyWO        types.String `tfsdk:"api_key_wo"`
	APIKeyWOVersion types.String `tfsdk:"api_key_wo_version"`
}

func (m *ConfigurationProfileJiraResource) Update(cp api.ConfigurationProfile) diag.Diagnostics {
	// fetch current project names to preserve declared order
	projectNames := typehelpers.ListItemsIdentifiers(m.Projects, "name")

	diags := m.ConfigurationProfileJiraDataSource.Update(cp)

	if projectNames != nil {
		projects, d := typehelpers.ListSortedEntries[JiraProject](m.Projects, "name", projectNames)
		m.Projects = projects
		diags.Append(d...)
	}

	return diags
}
