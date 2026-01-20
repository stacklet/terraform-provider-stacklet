// Copyright (c) 2026 - Stacklet, Inc.

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
)

var _ datasource.DataSource = &apiKeyDataSource{}

func newAPIKeyDataSource() datasource.DataSource {
	return &apiKeyDataSource{}
}

type apiKeyDataSource struct {
	apiDataSource
}

func (d *apiKeyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_key"
}

func (d *apiKeyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve information about an API key by its identity.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the API key.",
				Computed:    true,
			},
			"identity": schema.StringAttribute{
				Description: "The identity of the API key to look up.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the API key.",
				Computed:    true,
			},
			"expires_at": schema.StringAttribute{
				Description: "The timestamp when the API key expires (RFC3339 format).",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
			"revoked_at": schema.StringAttribute{
				Description: "The timestamp when the API key was revoked (RFC3339 format).",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
		},
	}
}

func (d *apiKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.APIKeyDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKey, err := d.api.APIKey.Read(ctx, data.Identity.ValueString())
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(data.Update(ctx, apiKey)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
