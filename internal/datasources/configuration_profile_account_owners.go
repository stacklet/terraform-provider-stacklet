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
	"github.com/stacklet/terraform-provider-stacklet/internal/modelupdate"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
)

var (
	_ datasource.DataSource = &configurationProfileAccountOwnersDataSource{}
)

func NewConfigurationProfileAccountOwnersDataSource() datasource.DataSource {
	return &configurationProfileAccountOwnersDataSource{}
}

type configurationProfileAccountOwnersDataSource struct {
	api *api.API
}

func (d *configurationProfileAccountOwnersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configuration_profile_account_owners"
}

func (d *configurationProfileAccountOwnersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve information about the Microsoft AccountOwners configuration profile.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the configuration profile.",
				Computed:    true,
			},
			"profile": schema.StringAttribute{
				Description: "The profile name.",
				Computed:    true,
			},
			"default": schema.ListNestedAttribute{
				Description: "List of default account owners.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"account": schema.StringAttribute{
							Description: "The account identifier.",
							Computed:    true,
						},
						"owners": schema.ListAttribute{
							Description: "List of owner addresses for this account.",
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
			},
			"org_domain": schema.StringAttribute{
				Description: "The organization domain to append to users for matching.",
				Computed:    true,
			},
			"org_domain_tag": schema.StringAttribute{
				Description: "The name of the tag to look up the organization domain.",
				Computed:    true,
			},
			"tags": schema.ListAttribute{
				Description: "List of tag names to look up the resource owner.",
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

func (d *configurationProfileAccountOwnersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if pd, err := providerdata.GetDataSourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		d.api = pd.API
	}
}

func (d *configurationProfileAccountOwnersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.ConfigurationProfileAccountOwnersDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, err := d.api.ConfigurationProfile.ReadAccountOwners(ctx)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	data.ID = types.StringValue(config.ID)
	data.Profile = types.StringValue(config.Profile)

	updater := modelupdate.NewConfigurationProfileUpdater(*config)
	defaultOwners, diags := updater.AccountOwnersDefault()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Default = defaultOwners

	data.OrgDomain = types.StringPointerValue(config.Record.AccountOwnersConfiguration.OrgDomain)
	data.OrgDomainTag = types.StringPointerValue(config.Record.AccountOwnersConfiguration.OrgDomainTag)
	data.Tags = tftypes.StringsList(config.Record.AccountOwnersConfiguration.Tags)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
