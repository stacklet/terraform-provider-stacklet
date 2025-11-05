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
	_ datasource.DataSource = &configurationProfileSymphonyDataSource{}
)

func NewConfigurationProfileSymphonyDataSource() datasource.DataSource {
	return &configurationProfileSymphonyDataSource{}
}

type configurationProfileSymphonyDataSource struct {
	api *api.API
}

func (d *configurationProfileSymphonyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configuration_profile_symphony"
}

func (d *configurationProfileSymphonyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve information about the Symphony configuration profile.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the configuration profile.",
				Computed:    true,
			},
			"profile": schema.StringAttribute{
				Description: "The profile name.",
				Computed:    true,
			},
			"agent_domain": schema.StringAttribute{
				Description: "The Symphony agent domain.",
				Computed:    true,
			},
			"service_account": schema.StringAttribute{
				Description: "The Symphony service account.",
				Computed:    true,
			},
			"private_key": schema.StringAttribute{
				Description: "The encrypted value for the account private key.",
				Computed:    true,
			},
		},
	}
}

func (d *configurationProfileSymphonyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if pd, err := providerdata.GetDataSourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		d.api = pd.API
	}
}

func (d *configurationProfileSymphonyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.ConfigurationProfileSymphonyDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, err := d.api.ConfigurationProfile.ReadSymphony(ctx)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(data.Update(*config)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
