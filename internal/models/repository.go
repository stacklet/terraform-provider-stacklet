// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

type RepositoryDataSource struct {
	ID                   types.String `tfsdk:"id"`
	UUID                 types.String `tfsdk:"uuid"`
	URL                  types.String `tfsdk:"url"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	WebhookURL           types.String `tfsdk:"webhook_url"`
	System               types.Bool   `tfsdk:"system"`
	AuthUser             types.String `tfsdk:"auth_user"`
	HasAuthToken         types.Bool   `tfsdk:"has_auth_token"`
	SSHPublicKey         types.String `tfsdk:"ssh_public_key"`
	HasSSHPrivateKey     types.Bool   `tfsdk:"has_ssh_private_key"`
	HasSSHPassphrase     types.Bool   `tfsdk:"has_ssh_passphrase"`
	RoleAssignmentTarget types.String `tfsdk:"role_assignment_target"`
}

func (m *RepositoryDataSource) Update(repo *api.Repository) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(repo.ID)
	m.UUID = types.StringValue(repo.UUID)
	m.URL = types.StringValue(repo.URL)
	m.Name = types.StringValue(repo.Name)
	m.Description = types.StringPointerValue(repo.Description)
	m.System = types.BoolValue(repo.System)
	m.WebhookURL = types.StringValue(repo.WebhookURL)
	m.AuthUser = types.StringPointerValue(repo.Auth.AuthUser)
	m.HasAuthToken = types.BoolValue(repo.Auth.HasAuthToken)
	m.SSHPublicKey = types.StringPointerValue(repo.Auth.SSHPublicKey)
	m.HasSSHPrivateKey = types.BoolValue(repo.Auth.HasSshPrivateKey)
	m.HasSSHPassphrase = types.BoolValue(repo.Auth.HasSshPassphrase)
	m.RoleAssignmentTarget = types.StringValue(repo.RoleAssignmentTarget)

	return diags
}

type RepositoryResource struct {
	RepositoryDataSource

	AuthTokenWO            types.String `tfsdk:"auth_token_wo"`
	AuthTokenWOVersion     types.String `tfsdk:"auth_token_wo_version"`
	SSHPrivateKeyWO        types.String `tfsdk:"ssh_private_key_wo"`
	SSHPrivateKeyWOVersion types.String `tfsdk:"ssh_private_key_wo_version"`
	SSHPassphraseWO        types.String `tfsdk:"ssh_passphrase_wo"`
	SSHPassphraseWOVersion types.String `tfsdk:"ssh_passphrase_wo_version"`
}
