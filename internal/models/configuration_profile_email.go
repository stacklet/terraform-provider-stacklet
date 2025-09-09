// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ConfigurationProfileEmailDataSource is the model for email configuration profile data sources.
type ConfigurationProfileEmailDataSource struct {
	ID        types.String `tfsdk:"id"`
	Profile   types.String `tfsdk:"profile"`
	From      types.String `tfsdk:"from"`
	SESRegion types.String `tfsdk:"ses_region"`
	SMTP      types.Object `tfsdk:"smtp"`
}

// SMTPDataSource is the model for SMTP configuration for data sources.
type SMTPDataSource struct {
	Server   types.String `tfsdk:"server"`
	Port     types.String `tfsdk:"port"`
	SSL      types.Bool   `tfsdk:"ssl"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (s SMTPDataSource) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"server":   types.StringType,
		"port":     types.StringType,
		"ssl":      types.BoolType,
		"username": types.StringType,
		"password": types.StringType,
	}
}

// SMTPResource is the model for SMTP configuration for resources.
type SMTPResource struct {
	SMTPDataSource

	PasswordWO        types.String `tfsdk:"password_wo"`
	PasswordWOVersion types.String `tfsdk:"password_wo_version"`
}

func (s SMTPResource) AttributeTypes() map[string]attr.Type {
	attrs := s.SMTPDataSource.AttributeTypes()
	attrs["password_wo"] = types.StringType
	attrs["password_wo_version"] = types.StringType
	return attrs
}

// ConfigurationProfileEmailResource is the model for email configuration profile resources.
type ConfigurationProfileEmailResource ConfigurationProfileEmailDataSource
