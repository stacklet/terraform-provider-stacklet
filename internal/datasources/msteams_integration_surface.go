// Copyright (c) 2025 - Stacklet, Inc.

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
)

var (
	_ datasource.DataSource = &msteamsIntegrationSurfaceDataSource{}
)

func newMSTeamsIntegrationSurfaceDataSource() datasource.DataSource {
	return &msteamsIntegrationSurfaceDataSource{}
}

type msteamsIntegrationSurfaceDataSource struct {
	apiDataSource
}

func (d *msteamsIntegrationSurfaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_msteams_integration_surface"
}

func (d *msteamsIntegrationSurfaceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve Microsoft Teams integration configuration details.",
		Attributes: map[string]schema.Attribute{
			"bot_endpoint": schema.StringAttribute{
				Description: "The bot endpoint URL for the MS Teams integration.",
				Computed:    true,
			},
			"oidc_client": schema.StringAttribute{
				Description: "The OIDC client identifier for the MS Teams integration.",
				Computed:    true,
			},
			"oidc_issuer": schema.StringAttribute{
				Description: "The OIDC issuer URL for the MS Teams integration.",
				Computed:    true,
			},
		},
	}
}

func (d *msteamsIntegrationSurfaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.MSTeamsIntegrationSurfaceDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	surface, err := d.api.System.MSTeamsIntegrationSurface(ctx)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(data.Update(ctx, surface)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
