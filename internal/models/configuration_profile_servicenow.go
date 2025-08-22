// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ConfigurationProfileServiceNowDataSource is the model for ServiceNow configuration profile data sources.
type ConfigurationProfileServiceNowDataSource struct {
	ID          types.String `tfsdk:"id"`
	Profile     types.String `tfsdk:"profile"`
	Endpoint    types.String `tfsdk:"endpoint"`
	Username    types.String `tfsdk:"username"`
	IssueType   types.String `tfsdk:"issue_type"`
	ClosedState types.String `tfsdk:"closed_state"`
}
