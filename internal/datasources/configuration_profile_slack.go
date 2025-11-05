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
)

var (
	_ datasource.DataSource = &configurationProfileSlackDataSource{}
)

func NewConfigurationProfileSlackDataSource() datasource.DataSource {
	return &configurationProfileSlackDataSource{}
}

type configurationProfileSlackDataSource struct {
	api *api.API
}

func (d *configurationProfileSlackDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configuration_profile_slack"
}

func (d *configurationProfileSlackDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve information about the Slack configuration profile.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the configuration profile.",
				Computed:    true,
			},

			"profile": schema.StringAttribute{
				Description: "The profile name.",
				Computed:    true,
			},
			"user_fields": schema.ListAttribute{
				Description: "Fields to use for identifying users for notification delivery.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"token": schema.StringAttribute{
				Description: "The encrypted value for the token.",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"webhook": schema.ListNestedBlock{
				Description: "Webhook configuration.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The webook name.",
							Computed:    true,
						},
						"url": schema.StringAttribute{
							Description: "The encrypted webhook URL.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *configurationProfileSlackDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if pd, err := providerdata.GetDataSourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		d.api = pd.API
	}
}

func (d *configurationProfileSlackDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.ConfigurationProfileSlackDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, err := d.api.ConfigurationProfile.ReadSlack(ctx)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(data.Update(*config)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
