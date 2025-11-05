// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
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

func (m *AccountDiscoveryGCPResource) Update(accountDiscovery *api.AccountDiscovery) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(accountDiscovery.ID)
	m.Name = types.StringValue(accountDiscovery.Name)
	m.Description = types.StringPointerValue(accountDiscovery.Description)
	m.Suspended = types.BoolValue(accountDiscovery.Schedule.Suspended)
	m.ClientEmail = types.StringValue(accountDiscovery.Config.GCPConfig.ClientEmail)
	m.ClientID = types.StringValue(accountDiscovery.Config.GCPConfig.ClientID)
	m.OrgID = types.StringValue(accountDiscovery.Config.GCPConfig.OrgID)
	m.RootFolderIDs = tftypes.StringsList(accountDiscovery.Config.GCPConfig.RootFolderIDs)
	m.ExcludeFolderIDs = tftypes.StringsList(accountDiscovery.Config.GCPConfig.ExcludeFolderIDs)
	m.ProjectID = types.StringValue(accountDiscovery.Config.GCPConfig.ProjectID)
	m.PrivateKeyID = types.StringValue(accountDiscovery.Config.GCPConfig.PrivateKeyID)

	return diags
}
