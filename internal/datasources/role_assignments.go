// Copyright (c) 2025 - Stacklet, Inc.

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
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
			"target": schema.SingleNestedAttribute{
				Description: "The target entity to query role assignments for.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "The type of target. Valid values: 'system', 'account-group', 'policy-collection', 'repository'.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("system", "account-group", "policy-collection", "repository"),
						},
					},
					"uuid": schema.StringAttribute{
						Description: "The UUID of the target. Required for all target types except 'system'.",
						Optional:    true,
					},
				},
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
						"principal": schema.SingleNestedAttribute{
							Description: "The principal (user or SSO group) that has been granted the role.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Description: "The type of principal. Either 'user' or 'sso-group'.",
									Computed:    true,
								},
								"id": schema.Int64Attribute{
									Description: "The numeric ID of the principal.",
									Computed:    true,
								},
							},
						},
						"target": schema.SingleNestedAttribute{
							Description: "The target entity that the role is assigned to.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Description: "The type of target.",
									Computed:    true,
								},
								"uuid": schema.StringAttribute{
									Description: "The UUID of the target (null for system targets).",
									Computed:    true,
								},
							},
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

	// Extract target from the model
	var target models.RoleAssignmentTarget
	resp.Diagnostics.Append(data.Target.As(ctx, &target, basetypes.ObjectAsOptions{})...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build API target filter
	apiTarget := &api.RoleAssignmentTarget{
		Type: target.Type.ValueString(),
		UUID: target.UUID.ValueStringPointer(),
	}

	// Query role assignments for the target
	assignments, err := d.api.RoleAssignment.List(ctx, apiTarget, nil)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	// Update the data model with results
	resp.Diagnostics.Append(data.Update(ctx, assignments)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
