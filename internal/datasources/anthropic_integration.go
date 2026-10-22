// Copyright Stacklet, Inc. 2025, 2026

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
)

var (
	_ datasource.DataSource              = &anthropicIntegrationDataSource{}
	_ datasource.DataSourceWithConfigure = &anthropicIntegrationDataSource{}
)

type anthropicIntegrationDataSource struct {
	apiDataSource
}

func (d *anthropicIntegrationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_anthropic_integration"
}

func (d *anthropicIntegrationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve an Anthropic integration.\n\n" +
			"The admin key is never returned by the API, and `configured` reports whether it holds one.\n",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the Anthropic integration.",
				Computed:    true,
			},
			"org_id": schema.StringAttribute{
				Description: "The ID of the Anthropic organization.",
				Required:    true,
			},
			"configured": schema.BoolAttribute{
				Description: "Whether the API holds an admin key.",
				Computed:    true,
			},
		},
	}
}

func (d *anthropicIntegrationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.AnthropicIntegrationDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integration, err := d.api.AnthropicIntegration.Read(ctx, data.OrgID.ValueString())
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(data.Update(integration)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
