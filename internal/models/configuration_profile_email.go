// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/typehelpers"
)

// ConfigurationProfileEmailDataSource is the model for email configuration profile data sources.
type ConfigurationProfileEmailDataSource struct {
	ID        types.String `tfsdk:"id"`
	Profile   types.String `tfsdk:"profile"`
	From      types.String `tfsdk:"from"`
	SESRegion types.String `tfsdk:"ses_region"`
	SMTP      types.Object `tfsdk:"smtp"`
}

func (m *ConfigurationProfileEmailDataSource) Update(ctx context.Context, cp api.ConfigurationProfile) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(cp.ID)
	m.Profile = types.StringValue(cp.Profile)
	m.From = types.StringValue(cp.Record.EmailConfiguration.FromEmail)
	m.SESRegion = types.StringPointerValue(cp.Record.EmailConfiguration.SESRegion)

	smtpConfig := cp.Record.EmailConfiguration.SMTP
	smtp, d := typehelpers.ObjectValue(
		ctx,
		smtpConfig,
		func() (*SMTPDataSource, diag.Diagnostics) {
			return &SMTPDataSource{
				Server:   types.StringValue(smtpConfig.Server),
				Port:     types.StringPointerValue(&smtpConfig.Port),
				SSL:      types.BoolPointerValue(smtpConfig.SSL),
				Username: types.StringPointerValue(smtpConfig.Username),
				Password: types.StringPointerValue(smtpConfig.Password),
			}, nil
		},
	)
	m.SMTP = smtp
	errors.AddAttributeDiags(&diags, d, "smtp")

	return diags
}

// ConfigurationProfileEmailResource is the model for email configuration profile resources.
type ConfigurationProfileEmailResource struct {
	ConfigurationProfileEmailDataSource
}

func (m *ConfigurationProfileEmailResource) Update(ctx context.Context, cp api.ConfigurationProfile) diag.Diagnostics {
	var diags diag.Diagnostics

	var origSMTPAttrs map[string]attr.Value

	if !m.SMTP.IsNull() {
		origSMTPAttrs = m.SMTP.Attributes()
	}

	diags.Append(m.ConfigurationProfileEmailDataSource.Update(ctx, cp)...)
	if diags.HasError() {
		return diags
	}

	if !m.SMTP.IsNull() {
		// Get original write-only fields if they existed, otherwise use null
		var origPasswordWO, origPasswordWOVersion types.String
		if origSMTPAttrs != nil {
			origPasswordWO, _ = origSMTPAttrs["password_wo"].(types.String)
			origPasswordWOVersion, _ = origSMTPAttrs["password_wo_version"].(types.String)
		}

		smtp, d := typehelpers.UpdatedObject(
			ctx,
			m.SMTP,
			map[string]attr.Value{
				"password_wo":         origPasswordWO,
				"password_wo_version": origPasswordWOVersion,
			},
		)
		errors.AddAttributeDiags(&diags, d, "smtp")
		m.SMTP = smtp
	}

	return diags
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
