// Copyright (c) 2025 - Stacklet, Inc.

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
)

var (
	_ datasource.DataSource = &accountGroupDataSource{}
)

func newAccountGroupDataSource() datasource.DataSource {
	return &accountGroupDataSource{}
}

type accountGroupDataSource struct {
	apiDataSource
}

func (d *accountGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account_group"
}

func (d *accountGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve an account group by UUID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the account group.",
				Computed:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "The UUID of the account group, alternative to the name.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the account group, alternative to the UUID.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the account group.",
				Computed:    true,
			},
			"dynamic_filter": schema.StringAttribute{
				Description: "Dynamic filter for accounts matching. Null means not dynamic, empty string matches all accounts.",
				Computed:    true,
			},
			"cloud_provider": schema.StringAttribute{
				Description: "The cloud provider for the account group.",
				Computed:    true,
			},
			"regions": schema.ListAttribute{
				Description: "The regions for the account group.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *accountGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.AccountGroupDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.UUID.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError("Invalid configuration", "Only one of UUID and name must be set")
		return
	}

	account_group, err := d.api.AccountGroup.Read(ctx, data.UUID.ValueString(), data.Name.ValueString())
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append((*models.AccountGroupResource)(&data).Update(account_group)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
