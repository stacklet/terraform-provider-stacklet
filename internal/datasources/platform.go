// Copyright (c) 2025 - Stacklet, Inc.

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
)

var (
	_ datasource.DataSource = &platformDataSource{}
)

func NewPlatformDataSource() datasource.DataSource {
	return &platformDataSource{}
}

type platformDataSource struct {
	api *api.API
}

func (d *platformDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_platform"
}

func (d *platformDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve information about the Stacklet platform.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID.",
				Computed:    true,
			},
			"external_id": schema.StringAttribute{
				Description: "The external ID for the deployment.",
				Computed:    true,
			},
			"execution_regions": schema.ListAttribute{
				Description: "List of regions for which execution is enabled.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"default_role": schema.StringAttribute{
				Description: "Default role for users.",
				Computed:    true,
			},
		},
	}
}

func (d *platformDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if pd, err := providerdata.GetDataSourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		d.api = pd.API
	}
}

func (d *platformDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.PlatformDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	platform, err := d.api.Platform.Read(ctx)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	data.ID = types.StringValue(platform.ID)
	data.ExternalID = tftypes.NullableString(platform.ExternalID)
	data.ExecutionRegions = tftypes.StringsList(platform.ExecutionRegions)
	data.DefaultRole = tftypes.NullableString(platform.DefaultRole)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
