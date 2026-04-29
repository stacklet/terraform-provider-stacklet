// Copyright Stacklet, Inc. 2025, 2026

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dsSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TerraformModule is the model for terraform modules definitions.
type TerraformModule struct {
	RepositoryURL types.String `tfsdk:"repository_url"`
	Source        types.String `tfsdk:"source"`
	VariablesJSON types.String `tfsdk:"variables_json"`
	Version       types.String `tfsdk:"version"`
}

func (c TerraformModule) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"repository_url": types.StringType,
		"source":         types.StringType,
		"variables_json": types.StringType,
		"version":        types.StringType,
	}
}

func (c TerraformModule) ResourceSchemaAttribute() resSchema.SingleNestedAttribute {
	return resSchema.SingleNestedAttribute{
		Description: "Terraform module configuration.",
		Computed:    true,
		Attributes: map[string]resSchema.Attribute{
			"repository_url": resSchema.StringAttribute{
				Description: "The repository URL.",
				Computed:    true,
			},
			"source": resSchema.StringAttribute{
				Description: "The module source.",
				Computed:    true,
			},
			"version": resSchema.StringAttribute{
				Description: "The module version.",
				Computed:    true,
			},
			"variables_json": resSchema.StringAttribute{
				Description: "The module variables as JSON.",
				Computed:    true,
			},
		},
	}
}

func (c TerraformModule) DataSourceSchemaAttribute() dsSchema.SingleNestedAttribute {
	return dsSchema.SingleNestedAttribute{
		Description: "Terraform module configuration.",
		Computed:    true,
		Attributes: map[string]dsSchema.Attribute{
			"repository_url": dsSchema.StringAttribute{
				Description: "The repository URL.",
				Computed:    true,
			},
			"source": dsSchema.StringAttribute{
				Description: "The module source.",
				Computed:    true,
			},
			"version": dsSchema.StringAttribute{
				Description: "The module version.",
				Computed:    true,
			},
			"variables_json": dsSchema.StringAttribute{
				Description: "The module variables as JSON.",
				Computed:    true,
			},
		},
	}
}
