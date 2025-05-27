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
	_ datasource.DataSource = &policyCollectionDataSource{}
)

func NewPolicyCollectionDataSource() datasource.DataSource {
	return &policyCollectionDataSource{}
}

type policyCollectionDataSource struct {
	api *api.API
}

func (d *policyCollectionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_collection"
}

func (d *policyCollectionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve a policy collection by UUID or name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the policy collection.",
				Computed:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "The UUID of the policy collection.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the policy collection.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the policy collection.",
				Computed:    true,
			},
			"cloud_provider": schema.StringAttribute{
				Description: "The cloud provider for the policy collection.",
				Computed:    true,
			},
			"auto_update": schema.BoolAttribute{
				Description: "Whether the policy collection automatically updates policy versions.",
				Computed:    true,
			},
			"system": schema.BoolAttribute{
				Description: "Whether this is a system policy collection.",
				Computed:    true,
			},
			"dynamic": schema.BoolAttribute{
				Description: "Whether this is a dynamic policy collection.",
				Computed:    true,
			},
			"repository_uuid": schema.StringAttribute{
				Description: "The UUID of the repository the collection is linked to, if dynamic.",
				Computed:    true,
			},
		},
	}
}

func (d *policyCollectionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if pd, err := providerdata.GetDataSourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		d.api = pd.API
	}
}

func (d *policyCollectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.PolicyCollectionDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyCollection, err := d.api.PolicyCollection.Read(ctx, data.UUID.ValueString(), data.Name.ValueString())
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	data.ID = types.StringValue(policyCollection.ID)
	data.UUID = types.StringValue(policyCollection.UUID)
	data.Name = types.StringValue(policyCollection.Name)
	data.Description = tftypes.NullableString(policyCollection.Description)
	data.CloudProvider = types.StringValue(string(policyCollection.Provider))
	data.AutoUpdate = types.BoolValue(policyCollection.AutoUpdate)
	data.System = types.BoolValue(policyCollection.System)
	data.Dynamic = types.BoolValue(policyCollection.IsDynamic)
	data.RepositoryUUID = tftypes.NullableString(policyCollection.RepositoryConfig.UUID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
