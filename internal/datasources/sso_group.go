// Copyright (c) 2025 - Stacklet, Inc.

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
)

var _ datasource.DataSource = &ssoGroupDataSource{}

func newSSOGroupDataSource() datasource.DataSource {
	return &ssoGroupDataSource{}
}

type ssoGroupDataSource struct {
	apiDataSource
}

func (d *ssoGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sso_group"
}

func (d *ssoGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve information about an SSO group by name. This data source provides the role_assignment_principal attribute needed for role assignments.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the SSO group.",
				Computed:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "The display name of the SSO group.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the SSO group.",
				Required:    true,
			},
			"role_assignment_principal": schema.StringAttribute{
				Description: "An opaque principal identifier for role assignments. Use this value when creating role assignments.",
				Computed:    true,
			},
		},
	}
}

func (d *ssoGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.SSOGroupDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ssoGroup, err := d.api.SSOGroup.Read(ctx, data.Name.ValueString())
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(data.Update(ctx, ssoGroup)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
