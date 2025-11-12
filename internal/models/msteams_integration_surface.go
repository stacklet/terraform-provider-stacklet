// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

// MSTeamsIntegrationSurfaceDataSource is the model for the MS Teams integration surface data source.
type MSTeamsIntegrationSurfaceDataSource struct {
	BotEndpoint types.String `tfsdk:"bot_endpoint"`
	OIDCClient  types.String `tfsdk:"oidc_client"`
	OIDCIssuer  types.String `tfsdk:"oidc_issuer"`
}

func (m *MSTeamsIntegrationSurfaceDataSource) Update(ctx context.Context, surface *api.MSTeamsIntegrationSurface) diag.Diagnostics {
	var diags diag.Diagnostics

	m.BotEndpoint = types.StringValue(surface.BotEndpoint)
	m.OIDCClient = types.StringValue(surface.OIDCClient)
	m.OIDCIssuer = types.StringValue(surface.OIDCIssuer)

	return diags
}
