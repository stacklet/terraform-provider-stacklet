// Copyright Stacklet, Inc. 2025, 2026

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &iacPolicyCollectionDataSource{}

// iacPolicyCollectionDataSource defines the data source implementation.
type iacPolicyCollectionDataSource struct {
	iacAPIDataSource
}

func (d *iacPolicyCollectionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iac_policy_collection"
}

func (d *iacPolicyCollectionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetch a Stacklet IaC policy collection, a set of IaC policies to evaluate.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the policy collection.",
				Required:    true,
			},
			"policies": schema.SetAttribute{
				Description: "The names of policies in the collection.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"is_default": schema.BoolAttribute{
				Description: "Whether the collection is a default, applied to projects that have no policy collections assigned.",
				Computed:    true,
			},
			"is_mirror": schema.BoolAttribute{
				Description: "Whether the collection mirrors the policies of a policy source.",
				Computed:    true,
			},
		},
	}
}

func (d *iacPolicyCollectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.IaCPolicyCollectionDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	collection, err := d.api.PolicyCollection.Read(ctx, data.Name.ValueString())
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(data.Update(ctx, collection)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
