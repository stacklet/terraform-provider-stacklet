// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AccountDiscoveryGCPResource is the model for GCP account discovery resources.
type AccountDiscoveryGCPResource struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Description           types.String `tfsdk:"description"`
	Suspended             types.Bool   `tfsdk:"suspended"`
	ClientEmail           types.String `tfsdk:"client_email"`
	ClientID              types.String `tfsdk:"client_id"`
	OrgID                 types.String `tfsdk:"org_id"`
	RootFolderIDs         types.List   `tfsdk:"root_folder_ids"`
	ExcludeFolderIDs      types.List   `tfsdk:"exclude_folder_ids"`
	ProjectID             types.String `tfsdk:"project_id"`
	PrivateKeyID          types.String `tfsdk:"private_key_id"`
	CredentialJSON        types.String `tfsdk:"credential_json_wo"`
	CredentialJSONVersion types.String `tfsdk:"credential_json_wo_version"`
}
