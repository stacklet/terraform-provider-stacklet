// Copyright Stacklet, Inc. 2025, 2026

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
)

var _ datasource.DataSource = &userGroupDataSource{}

type userGroupDataSource struct {
	apiDataSource
}

func (d *userGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_group"
}

func (d *userGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve information about a user group by name. This data source provides the role_assignment_principal and role_assignment_target attributes needed for role assignments.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the user group.",
				Computed:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "The UUID of the user group.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the user group.",
				Required:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "The display name of the user group.",
				Computed:    true,
			},
			"role_assignment_principal": schema.StringAttribute{
				Description: "An opaque principal identifier for role assignments. Use this value when granting roles to the members of this group.",
				Computed:    true,
			},
			"role_assignment_target": schema.StringAttribute{
				Description: "An opaque target identifier for role assignments. Use this value when granting roles on this group.",
				Computed:    true,
			},
		},
	}
}

func (d *userGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.UserGroupDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userGroup, err := d.api.UserGroup.ReadByName(ctx, data.Name.ValueString())
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(data.Update(userGroup)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
