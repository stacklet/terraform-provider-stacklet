// Copyright (c) 2025 - Stacklet, Inc.

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
)

var (
	_ datasource.DataSource = &configurationProfileServiceNowDataSource{}
)

func NewConfigurationProfileServiceNowDataSource() datasource.DataSource {
	return &configurationProfileServiceNowDataSource{}
}

type configurationProfileServiceNowDataSource struct {
	api *api.API
}

func (d *configurationProfileServiceNowDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configuration_profile_servicenow"
}

func (d *configurationProfileServiceNowDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve information about the ServiceNow configuration profile.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the configuration profile.",
				Computed:    true,
			},
			"profile": schema.StringAttribute{
				Description: "The profile name.",
				Computed:    true,
			},
			"endpoint": schema.StringAttribute{
				Description: "The ServiceNow instance endpoint.",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "The ServiceNow instance authentication username.",
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description: "The encrypted value for the ServiceNow instance authentication password.",
				Computed:    true,
			},
			"issue_type": schema.StringAttribute{
				Description: "The type of issue to use for tickets.",
				Computed:    true,
			},
			"closed_state": schema.StringAttribute{
				Description: "The state for closed tickets.",
				Computed:    true,
			},
		},
	}
}

func (d *configurationProfileServiceNowDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if pd, err := providerdata.GetDataSourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		d.api = pd.API
	}
}

func (d *configurationProfileServiceNowDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.ConfigurationProfileServiceNowDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, err := d.api.ConfigurationProfile.ReadServiceNow(ctx)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(data.Update(*config)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
