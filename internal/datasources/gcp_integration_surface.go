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
	_ datasource.DataSource = &gcpIntegrationSurfaceDataSource{}
)

func newGCPIntegrationSurfaceDataSource() datasource.DataSource {
	return &gcpIntegrationSurfaceDataSource{}
}

type gcpIntegrationSurfaceDataSource struct {
	apiDataSource
}

func (d *gcpIntegrationSurfaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gcp_integration_surface"
}

func (d *gcpIntegrationSurfaceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve GCP integration configuration details.",
		Attributes: map[string]schema.Attribute{
			"trust_aws": schema.SingleNestedAttribute{
				Description: "AWS trust configuration for GCP integration.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"account_id": schema.StringAttribute{
						Description: "The Stacklet deployment AWS account ID that GCP should integrate with.",
						Computed:    true,
					},
					"assetdb_role_name": schema.StringAttribute{
						Description: "The name of the IAM role used for AssetDB access.",
						Computed:    true,
					},
					"cost_query_role_name": schema.StringAttribute{
						Description: "The name of the IAM role used for cost queries access.",
						Computed:    true,
					},
					"execution_role_name": schema.StringAttribute{
						Description: "The name of the IAM role used for execution.",
						Computed:    true,
					},
					"platform_role_name": schema.StringAttribute{
						Description: "The name of the IAM role used for platform access.",
						Computed:    true,
					},
				},
			},
			"aws_relay": schema.SingleNestedAttribute{
				Description: "AWS relay configuration for GCP integration.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"bus_arn": schema.StringAttribute{
						Description: "The ARN of the AWS EventBridge bus to relay events to.",
						Computed:    true,
					},
					"role_arn": schema.StringAttribute{
						Description: "The ARN of the AWS IAM role used for the relay.",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d *gcpIntegrationSurfaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.GCPIntegrationSurfaceDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	surface, err := d.api.System.GCPIntegrationSurface(ctx)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(data.Update(ctx, surface)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
