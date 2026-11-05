// Copyright Stacklet, Inc. 2025, 2026

package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/stacklet/terraform-provider-stacklet/internal/api/iac"
	"github.com/stacklet/terraform-provider-stacklet/internal/typehelpers"
)

// IaCPolicySourceBase holds attributes common to the policy source resource and data source.
type IaCPolicySourceBase struct {
	Name       types.String `tfsdk:"name"`
	URL        types.String `tfsdk:"url"`
	Branch     types.String `tfsdk:"branch"`
	Paths      types.Set    `tfsdk:"paths"`
	Extensions types.Set    `tfsdk:"extensions"`
	AutoScan   types.Bool   `tfsdk:"auto_scan"`
	WebhookURL types.String `tfsdk:"webhook_url"`
}

func (m *IaCPolicySourceBase) Update(ctx context.Context, source *iac.PolicySource) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Name = types.StringValue(source.Name)
	m.URL = types.StringValue(source.URL)
	m.Branch = types.StringPointerValue(source.Branch)
	m.AutoScan = types.BoolValue(source.AutoScan)
	m.WebhookURL = types.StringValue(source.WebhookURL)

	var d diag.Diagnostics
	if source.Paths == nil {
		m.Paths = types.SetNull(types.StringType)
	} else {
		m.Paths, d = types.SetValueFrom(ctx, types.StringType, source.Paths)
		diags.Append(d...)
	}
	m.Extensions, d = types.SetValueFrom(ctx, types.StringType, source.Extensions)
	diags.Append(d...)

	return diags
}

// splitPolicySourceAuth returns the policy source auth as either token or SSH auth.
func splitPolicySourceAuth(auth *iac.PolicySourceAuth) (tokenAuth, sshAuth *iac.PolicySourceAuth) {
	switch {
	case auth == nil:
		return nil, nil
	case auth.Token != nil:
		return auth, nil
	default:
		return nil, auth
	}
}

type IaCPolicySourceDataSource struct {
	IaCPolicySourceBase

	TokenAuth types.Object `tfsdk:"token_auth"`
	SSHAuth   types.Object `tfsdk:"ssh_auth"`
}

func (m *IaCPolicySourceDataSource) Update(ctx context.Context, source *iac.PolicySource) diag.Diagnostics {
	var diags diag.Diagnostics

	diags.Append(m.IaCPolicySourceBase.Update(ctx, source)...)

	tokenAuth, sshAuth := splitPolicySourceAuth(source.Auth)

	var d diag.Diagnostics
	m.TokenAuth, d = typehelpers.ObjectValue(
		ctx,
		tokenAuth,
		func() (*IaCPolicySourceTokenAuthDataSource, diag.Diagnostics) {
			return &IaCPolicySourceTokenAuthDataSource{
				User: types.StringValue(tokenAuth.Username),
			}, nil
		},
	)
	diags.Append(d...)
	m.SSHAuth, d = typehelpers.ObjectValue(
		ctx,
		sshAuth,
		func() (*IaCPolicySourceSSHAuthDataSource, diag.Diagnostics) {
			return &IaCPolicySourceSSHAuthDataSource{
				User:          types.StringValue(sshAuth.Username),
				HasPassphrase: types.BoolValue(sshAuth.SSHPassphrase != nil),
			}, nil
		},
	)
	diags.Append(d...)

	return diags
}

type IaCPolicySourceTokenAuthDataSource struct {
	User types.String `tfsdk:"user"`
}

func (a IaCPolicySourceTokenAuthDataSource) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"user": types.StringType,
	}
}

type IaCPolicySourceSSHAuthDataSource struct {
	User          types.String `tfsdk:"user"`
	HasPassphrase types.Bool   `tfsdk:"has_passphrase"`
}

func (a IaCPolicySourceSSHAuthDataSource) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"user":           types.StringType,
		"has_passphrase": types.BoolType,
	}
}

type IaCPolicySourceResource struct {
	IaCPolicySourceBase

	TokenAuth types.Object `tfsdk:"token_auth"`
	SSHAuth   types.Object `tfsdk:"ssh_auth"`
}

// Update updates the model from the API response.
//
// Secrets are not returned by the API, so _wo_version attributes are
// preserved from the current model, as long as the corresponding secret is
// still set.
func (m *IaCPolicySourceResource) Update(ctx context.Context, source *iac.PolicySource) diag.Diagnostics {
	var diags diag.Diagnostics

	diags.Append(m.IaCPolicySourceBase.Update(ctx, source)...)

	tokenAuth, sshAuth := splitPolicySourceAuth(source.Auth)

	var currentToken IaCPolicySourceTokenAuth
	if !m.TokenAuth.IsNull() && !m.TokenAuth.IsUnknown() {
		diags.Append(m.TokenAuth.As(ctx, &currentToken, basetypes.ObjectAsOptions{})...)
	}
	var currentSSH IaCPolicySourceSSHAuth
	if !m.SSHAuth.IsNull() && !m.SSHAuth.IsUnknown() {
		diags.Append(m.SSHAuth.As(ctx, &currentSSH, basetypes.ObjectAsOptions{})...)
	}
	if diags.HasError() {
		return diags
	}

	var d diag.Diagnostics
	m.TokenAuth, d = typehelpers.ObjectValue(
		ctx,
		tokenAuth,
		func() (*IaCPolicySourceTokenAuth, diag.Diagnostics) {
			return &IaCPolicySourceTokenAuth{
				User:           types.StringValue(tokenAuth.Username),
				TokenWO:        types.StringNull(),
				TokenWOVersion: currentToken.TokenWOVersion,
			}, nil
		},
	)
	diags.Append(d...)
	m.SSHAuth, d = typehelpers.ObjectValue(
		ctx,
		sshAuth,
		func() (*IaCPolicySourceSSHAuth, diag.Diagnostics) {
			passphraseVersion := types.StringNull()
			if sshAuth.SSHPassphrase != nil {
				passphraseVersion = currentSSH.PassphraseWOVersion
			}
			return &IaCPolicySourceSSHAuth{
				User:                types.StringValue(sshAuth.Username),
				PrivateKeyWO:        types.StringNull(),
				PrivateKeyWOVersion: currentSSH.PrivateKeyWOVersion,
				PassphraseWO:        types.StringNull(),
				PassphraseWOVersion: passphraseVersion,
			}, nil
		},
	)
	diags.Append(d...)

	return diags
}

type IaCPolicySourceTokenAuth struct {
	User           types.String `tfsdk:"user"`
	TokenWO        types.String `tfsdk:"token_wo"`
	TokenWOVersion types.String `tfsdk:"token_wo_version"`
}

func (a IaCPolicySourceTokenAuth) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"user":             types.StringType,
		"token_wo":         types.StringType,
		"token_wo_version": types.StringType,
	}
}

type IaCPolicySourceSSHAuth struct {
	User                types.String `tfsdk:"user"`
	PrivateKeyWO        types.String `tfsdk:"private_key_wo"`
	PrivateKeyWOVersion types.String `tfsdk:"private_key_wo_version"`
	PassphraseWO        types.String `tfsdk:"passphrase_wo"`
	PassphraseWOVersion types.String `tfsdk:"passphrase_wo_version"`
}

func (a IaCPolicySourceSSHAuth) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"user":                   types.StringType,
		"private_key_wo":         types.StringType,
		"private_key_wo_version": types.StringType,
		"passphrase_wo":          types.StringType,
		"passphrase_wo_version":  types.StringType,
	}
}
