// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

// APIKeyDataSource is the model for API key data sources.
type APIKeyDataSource struct {
	ID          types.String      `tfsdk:"id"`
	Identity    types.String      `tfsdk:"identity"`
	Description types.String      `tfsdk:"description"`
	ExpiresAt   timetypes.RFC3339 `tfsdk:"expires_at"`
	RevokedAt   timetypes.RFC3339 `tfsdk:"revoked_at"`
}

func (m *APIKeyDataSource) Update(ctx context.Context, apiKey *api.APIKey) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(apiKey.ID)
	m.Identity = types.StringValue(apiKey.Identity)
	m.Description = types.StringValue(apiKey.Description)

	var d diag.Diagnostics
	m.ExpiresAt, d = timetypes.NewRFC3339PointerValue(apiKey.ExpiresAt)
	diags.Append(d...)
	m.RevokedAt, d = timetypes.NewRFC3339PointerValue(apiKey.RevokedAt)
	diags.Append(d...)

	return diags
}

// APIKeyResource is the model for API key resources.
type APIKeyResource struct {
	APIKeyDataSource

	Secret types.String `tfsdk:"secret"`
}

func (m *APIKeyResource) Update(ctx context.Context, apiKey *api.APIKey, secret *string) diag.Diagnostics {
	diag := m.APIKeyDataSource.Update(ctx, apiKey)

	if secret != nil {
		m.Secret = types.StringPointerValue(secret)
	}
	return diag
}
