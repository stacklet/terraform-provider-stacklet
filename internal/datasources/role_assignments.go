// Copyright (c) 2025 - Stacklet, Inc.

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
)

var (
	_ datasource.DataSource = &roleAssignmentsDataSource{}
)

func newRoleAssignmentsDataSource() datasource.DataSource {
	return &roleAssignmentsDataSource{}
}

type roleAssignmentsDataSource struct {
	apiDataSource
}

func (d *roleAssignmentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role_assignments"
}

func (d *roleAssignmentsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve role assignments for a specific target. This data source allows you to query which principals (users or SSO groups) have been granted roles on a particular target (system, account group, policy collection, or repository).",
		Attributes: map[string]schema.Attribute{
			"target": schema.StringAttribute{
				Description: "An opaque target identifier to query role assignments for. Use the 'role_assignment_target' attribute from resource outputs.",
				Required:    true,
			},
			"assignments": schema.ListNestedAttribute{
				Description: "The list of role assignments for the target.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier of the role assignment.",
							Computed:    true,
						},
						"role_name": schema.StringAttribute{
							Description: "The name of the role assigned.",
							Computed:    true,
						},
						"principal": schema.StringAttribute{
							Description: "An opaque principal identifier.",
							Computed:    true,
						},
						"target": schema.StringAttribute{
							Description: "An opaque target identifier.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *roleAssignmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.RoleAssignmentsDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Query role assignments for the target
	// We pass the opaque target string directly to the API for filtering
	targetStr := data.Target.ValueString()
	assignments, err := d.api.RoleAssignment.List(ctx, &targetStr, nil)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	// Update the data model with results
	resp.Diagnostics.Append(data.Update(ctx, assignments)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
