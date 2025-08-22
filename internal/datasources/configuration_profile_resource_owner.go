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
	_ datasource.DataSource = &configurationProfileResourceOwnerDataSource{}
)

func NewConfigurationProfileResourceOwnerDataSource() datasource.DataSource {
	return &configurationProfileResourceOwnerDataSource{}
}

type configurationProfileResourceOwnerDataSource struct {
	api *api.API
}

func (d *configurationProfileResourceOwnerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configuration_profile_resource_owner"
}

func (d *configurationProfileResourceOwnerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve information about the Microsoft ResourceOwner configuration profile.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the configuration profile.",
				Computed:    true,
			},
			"profile": schema.StringAttribute{
				Description: "The profile name.",
				Computed:    true,
			},
			"default": schema.ListAttribute{
				Description: "List of fallback notification addresses.",
				ElementType: types.StringType,
				Computed:    true,
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

func (d *configurationProfileResourceOwnerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if pd, err := providerdata.GetDataSourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		d.api = pd.API
	}
}

func (d *configurationProfileResourceOwnerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.ConfigurationProfileResourceOwnerDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, err := d.api.ConfigurationProfile.ReadResourceOwner(ctx)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	data.ID = types.StringValue(config.ID)
	data.Profile = types.StringValue(config.Profile)
	data.Default = tftypes.StringsList(config.Record.ResourceOwnerConfiguration.Default)
	data.OrgDomain = types.StringPointerValue(config.Record.ResourceOwnerConfiguration.OrgDomain)
	data.OrgDomainTag = types.StringPointerValue(config.Record.ResourceOwnerConfiguration.OrgDomainTag)
	data.Tags = tftypes.StringsList(config.Record.ResourceOwnerConfiguration.Tags)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
