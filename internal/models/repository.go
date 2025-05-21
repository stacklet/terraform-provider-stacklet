package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RepositoryDataSource struct {
	ID               types.String `tfsdk:"id"`
	UUID             types.String `tfsdk:"uuid"`
	URL              types.String `tfsdk:"url"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	WebhookURL       types.String `tfsdk:"webhook_url"`
	System           types.Bool   `tfsdk:"system"`
	AuthUser         types.String `tfsdk:"auth_user"`
	HasAuthToken     types.Bool   `tfsdk:"has_auth_token"`
	SSHPublicKey     types.String `tfsdk:"ssh_public_key"`
	HasSSHPrivateKey types.Bool   `tfsdk:"has_ssh_private_key"`
	HasSSHPassphrase types.Bool   `tfsdk:"has_ssh_passphrase"`
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
