// Copyright (c) 2025 - Stacklet, Inc.

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

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
			"dynamic_config": schema.SingleNestedAttribute{
				Description: "Configuration for dynamic behavior.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"repository_uuid": schema.StringAttribute{
						Description: "The UUID of the repository the collection is linked to.",
						Computed:    true,
					},
					"namespace": schema.StringAttribute{
						Description: "The namespace for policies from the repository.",
						Computed:    true,
					},
					"branch_name": schema.StringAttribute{
						Description: "The repository branch from which policies are imported.",
						Computed:    true,
					},
					"policy_directories": schema.ListAttribute{
						Description: "Optional list of subdirectory to limit the scan to.",
						Computed:    true,
						ElementType: types.StringType,
					},
					"policy_file_suffixes": schema.ListAttribute{
						Description: "Optional list of suffixes for policy files to limit the scan to.",
						Computed:    true,
						ElementType: types.StringType,
					},
				},
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

	attrTypes := models.PolicyCollectionDynamicConfig{}.AttributeTypes()
	var config basetypes.ObjectValue
	var diags diag.Diagnostics
	if policyCollection.RepositoryView == nil {
		config = basetypes.NewObjectNull(attrTypes)
	} else {
		config, diags = basetypes.NewObjectValueFrom(
			ctx,
			attrTypes,
			models.PolicyCollectionDynamicConfig{
				RepositoryUUID:     types.StringValue(*policyCollection.RepositoryConfig.UUID),
				Namespace:          types.StringValue(policyCollection.RepositoryView.Namespace),
				BranchName:         types.StringValue(policyCollection.RepositoryView.BranchName),
				PolicyDirectories:  tftypes.StringsList(policyCollection.RepositoryView.PolicyDirectories),
				PolicyFileSuffixes: tftypes.StringsList(policyCollection.RepositoryView.PolicyFileSuffix),
			},
		)
	}
	data.DynamicConfig = config

	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
