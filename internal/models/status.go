// Copyright Stacklet, Inc. 2025, 2026

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dsSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	resSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/typehelpers"
)

const (
	statusInfoDescription = "Status of the background operations that put this configuration in force, and their rolled-up summary."

	statusSummaryDescription = "Rolled up across the details: `ERROR` if every operation that reached an outcome failed, `PARTIAL_SUCCESS` if successes mix with failures or with operations still outstanding, `SUCCESS` if every operation that reached an outcome succeeded, and `SKIPPED` when there is no outcome anywhere yet. An operation still waiting on its first outcome keeps the summary off `SUCCESS` but never counts as a failure."

	statusDetailsDescription = "One entry per operation. An operation with nothing to report — never run, and nothing underway — has no entry."

	statusDetailStatusDescription = "Outcome of the most recent completed attempt; null when none has completed. Not a state: `running_since` says whether an attempt is underway."

	statusDetailAtDescription = "When that attempt finished; null whenever `status` is, and when `status` is `UNAVAILABLE`."

	statusDetailMessageDescription = "A human-readable summary of that attempt, typically why it failed. May be set alongside a null `status`, to say what is being waited on."

	statusDetailRunningSinceDescription = "Set while an attempt is underway, whether or not one has completed before; the other fields go on describing the last completed attempt."
)

// StatusInfo is the model for the status of an entity's background operations.
type StatusInfo struct {
	Status  types.String `tfsdk:"status"`
	Details types.List   `tfsdk:"details"`
}

func (s StatusInfo) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"status":  types.StringType,
		"details": types.ListType{ElemType: types.ObjectType{AttrTypes: StatusDetail{}.AttributeTypes()}},
	}
}

// StatusDetail is the model for the outcome of a single named operation.
type StatusDetail struct {
	Component    types.String `tfsdk:"component"`
	Status       types.String `tfsdk:"status"`
	At           types.String `tfsdk:"at"`
	Message      types.String `tfsdk:"message"`
	RunningSince types.String `tfsdk:"running_since"`
}

func (s StatusDetail) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"component":     types.StringType,
		"status":        types.StringType,
		"at":            types.StringType,
		"message":       types.StringType,
		"running_since": types.StringType,
	}
}

// NewStatusInfo returns the object value for a status returned by the API.
func NewStatusInfo(info api.StatusInfo) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	details, d := typehelpers.ObjectList[StatusDetail](
		info.Details,
		func(detail api.StatusDetail) (map[string]attr.Value, diag.Diagnostics) {
			return map[string]attr.Value{
				"component":     types.StringValue(detail.Component),
				"status":        types.StringPointerValue(detail.Status),
				"at":            types.StringPointerValue(detail.At),
				"message":       types.StringPointerValue(detail.Message),
				"running_since": types.StringPointerValue(detail.RunningSince),
			}, nil
		},
	)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(StatusInfo{}.AttributeTypes()), diags
	}

	return types.ObjectValue(
		StatusInfo{}.AttributeTypes(),
		map[string]attr.Value{
			"status":  types.StringValue(info.Status),
			"details": details,
		},
	)
}

func (s StatusInfo) ResourceSchemaAttribute(detailDescription string) resSchema.SingleNestedAttribute {
	return resSchema.SingleNestedAttribute{
		Description: statusInfoDescription,
		Computed:    true,
		Attributes: map[string]resSchema.Attribute{
			"status": resSchema.StringAttribute{
				Description: statusSummaryDescription,
				Computed:    true,
			},
			"details": resSchema.ListNestedAttribute{
				Description: statusDetailsDescription,
				Computed:    true,
				NestedObject: resSchema.NestedAttributeObject{
					Attributes: map[string]resSchema.Attribute{
						"component": resSchema.StringAttribute{
							Description: detailDescription,
							Computed:    true,
						},
						"status": resSchema.StringAttribute{
							Description: statusDetailStatusDescription,
							Computed:    true,
						},
						"at": resSchema.StringAttribute{
							Description: statusDetailAtDescription,
							Computed:    true,
						},
						"message": resSchema.StringAttribute{
							Description: statusDetailMessageDescription,
							Computed:    true,
						},
						"running_since": resSchema.StringAttribute{
							Description: statusDetailRunningSinceDescription,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (s StatusInfo) DataSourceSchemaAttribute(detailDescription string) dsSchema.SingleNestedAttribute {
	return dsSchema.SingleNestedAttribute{
		Description: statusInfoDescription,
		Computed:    true,
		Attributes: map[string]dsSchema.Attribute{
			"status": dsSchema.StringAttribute{
				Description: statusSummaryDescription,
				Computed:    true,
			},
			"details": dsSchema.ListNestedAttribute{
				Description: statusDetailsDescription,
				Computed:    true,
				NestedObject: dsSchema.NestedAttributeObject{
					Attributes: map[string]dsSchema.Attribute{
						"component": dsSchema.StringAttribute{
							Description: detailDescription,
							Computed:    true,
						},
						"status": dsSchema.StringAttribute{
							Description: statusDetailStatusDescription,
							Computed:    true,
						},
						"at": dsSchema.StringAttribute{
							Description: statusDetailAtDescription,
							Computed:    true,
						},
						"message": dsSchema.StringAttribute{
							Description: statusDetailMessageDescription,
							Computed:    true,
						},
						"running_since": dsSchema.StringAttribute{
							Description: statusDetailRunningSinceDescription,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
