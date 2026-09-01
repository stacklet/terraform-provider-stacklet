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
	_ datasource.DataSource              = &datadogIntegrationDataSource{}
	_ datasource.DataSourceWithConfigure = &datadogIntegrationDataSource{}
)

type datadogIntegrationDataSource struct {
	apiDataSource
}

func (d *datadogIntegrationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_datadog_integration"
}

func (d *datadogIntegrationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve the Datadog integration.\n\n" +
			"The integration applies to the whole deployment, so this takes no arguments. " +
			"The credentials are never returned by the API, and `configured` reports whether it holds them.\n",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "A fixed identifier for the integration, which is global and has no ID of its own.",
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether policies may use the `datadog-metrics` filter.",
				Computed:    true,
			},
			"site": schema.StringAttribute{
				Description: "The Datadog site the credentials belong to, as a bare hostname, e.g. `datadoghq.com`.",
				Computed:    true,
			},
			"configured": schema.BoolAttribute{
				Description: "Whether the API holds credentials.",
				Computed:    true,
			},
			"status": models.StatusInfo{}.DataSourceSchemaAttribute(
				"What the entry describes: `<region>-secret`, e.g. `us-east-1-secret`, is whether Secrets Manager holds the credentials in that region, which it replicates on its own schedule; `<region>-execution` is whether the region has recorded where to find them, so a policy running there can use them. Every region the configuration is meant to reach has both from the moment that is decided, so the list doesn't grow as propagation lands.",
			),
		},
	}
}

func (d *datadogIntegrationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.DatadogIntegrationDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integration, err := d.api.DatadogIntegration.Read(ctx)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(data.Update(integration)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
