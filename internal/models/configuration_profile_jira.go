// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ConfigurationProfileJiraDataSource is the model for Jira configuration profile data sources.
type ConfigurationProfileJiraDataSource struct {
	ID       types.String `tfsdk:"id"`
	Profile  types.String `tfsdk:"profile"`
	URL      types.String `tfsdk:"url"`
	Projects types.List   `tfsdk:"project"`
	User     types.String `tfsdk:"user"`
	APIKey   types.String `tfsdk:"api_key"`
}

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

// ConfigurationProfileJiraResource is the model for Jira configuration profile resources.
type ConfigurationProfileJiraResource struct {
	ConfigurationProfileJiraDataSource

	APIKeyWO        types.String `tfsdk:"api_key_wo"`
	APIKeyWOVersion types.String `tfsdk:"api_key_wo_version"`
}
