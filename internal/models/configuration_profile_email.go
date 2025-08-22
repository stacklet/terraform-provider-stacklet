// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ConfigurationProfileEmailDataSource is the model for email configuration profile data sources.
type ConfigurationProfileEmailDataSource struct {
	ID      types.String `tfsdk:"id"`
	Profile types.String `tfsdk:"profile"`
	From    types.String `tfsdk:"from"`
	SMTP    types.Object `tfsdk:"smtp"`
}

// SMTP is the model for SMTP configuration.
type SMTP struct {
	Server   types.String `tfsdk:"server"`
	Port     types.String `tfsdk:"port"`
	SSL      types.Bool   `tfsdk:"ssl"`
	Username types.String `tfsdk:"username"`
}

func (s SMTP) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"server":   types.StringType,
		"port":     types.StringType,
		"ssl":      types.BoolType,
		"username": types.StringType,
	}
}
