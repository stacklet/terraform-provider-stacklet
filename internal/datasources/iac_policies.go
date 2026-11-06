// Copyright Stacklet, Inc. 2025, 2026

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api/iac"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &iacPoliciesDataSource{}

// iacPoliciesDataSource defines the data source implementation.
type iacPoliciesDataSource struct {
	iacAPIDataSource
}

func (d *iacPoliciesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iac_policies"
}

func (d *iacPoliciesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetch Stacklet IaC policies from a policy source or a policy collection.",
		Attributes: map[string]schema.Attribute{
			"source": schema.StringAttribute{
				Description: "The name of the policy source to fetch policies from. Exactly one of `source` or `collection` must be set.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("collection")),
				},
			},
			"collection": schema.StringAttribute{
				Description: "The name of the policy collection to fetch policies from. Exactly one of `source` or `collection` must be set.",
				Optional:    true,
			},
			"policies": schema.ListNestedAttribute{
				Description: "The policies, sorted by name.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the policy.",
							Computed:    true,
						},
						"source": schema.StringAttribute{
							Description: "The name of the policy source the policy comes from.",
							Computed:    true,
						},
						"resources": schema.ListAttribute{
							Description: "The resource types the policy applies to.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"severity": schema.StringAttribute{
							Description: "The severity of the policy.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *iacPoliciesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.IaCPoliciesDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var policies []iac.Policy
	var err error
	if !data.Source.IsNull() {
		policies, err = d.api.Policy.ListBySource(ctx, data.Source.ValueString())
	} else {
		policies, err = d.api.Policy.ListByCollection(ctx, data.Collection.ValueString())
	}
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(data.Update(ctx, policies)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
