// Copyright Stacklet, Inc. 2025, 2026

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/typehelpers"
)

// SAMLProviderDataSource is the model for SAML provider data sources.
type SAMLProviderDataSource struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	DisplayName   types.String `tfsdk:"display_name"`
	MetadataURL   types.String `tfsdk:"metadata_url"`
	MetadataXML   types.String `tfsdk:"metadata_xml"`
	EnableSignout types.Bool   `tfsdk:"enable_signout"`
	IDPAlias      types.String `tfsdk:"idp_alias"`
}

func (m *SAMLProviderDataSource) Update(samlProvider *api.SAMLProvider) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = typehelpers.GraphQLIDValue(samlProvider.ID)
	m.Name = types.StringValue(samlProvider.Name)
	m.DisplayName = types.StringPointerValue(samlProvider.DisplayName)
	m.MetadataURL = types.StringPointerValue(samlProvider.MetadataURL)
	m.MetadataXML = types.StringPointerValue(samlProvider.MetadataXML)
	m.EnableSignout = types.BoolValue(samlProvider.EnableSignout)
	m.IDPAlias = types.StringPointerValue(samlProvider.IDPAlias)

	return diags
}

// SAMLProviderResource is the model for SAML provider resources.
type SAMLProviderResource struct {
	SAMLProviderDataSource
}
