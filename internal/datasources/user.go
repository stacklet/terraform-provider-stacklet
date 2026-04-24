// Copyright (c) 2025 - Stacklet, Inc.

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
)

var _ datasource.DataSource = &userDataSource{}

func newUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

type userDataSource struct {
	apiDataSource
}

func (d *userDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *userDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve information about a user by username. This data source provides the role_assignment_principal attribute needed for role assignments.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the user.",
				Computed:    true,
			},
			"active": schema.BoolAttribute{
				Description: "Whether the user is active in the system.",
				Computed:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "The display name of the user.",
				Computed:    true,
			},
			"email": schema.StringAttribute{
				Description: "The email address of the user.",
				Computed:    true,
			},
			"key": schema.Int64Attribute{
				Description: "The numeric key identifier of the user.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the user.",
				Computed:    true,
			},
			"role_assignment_principal": schema.StringAttribute{
				Description: "An opaque principal identifier for role assignments. Use this value when creating role assignments.",
				Computed:    true,
			},
			"sso_user": schema.BoolAttribute{
				Description: "Whether the user is an SSO user.",
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "The username of the user to look up.",
				Required:    true,
			},
		},
	}
}

func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.UserDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Query user by username
	user, err := d.api.User.Read(ctx, data.Username.ValueString())
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(data.Update(ctx, user)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
