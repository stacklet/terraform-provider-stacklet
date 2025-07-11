// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// PlatformDataSource is the model for the platform data source.
type PlatformDataSource struct {
	ID                       types.String `tfsdk:"id"`
	ExternalID               types.String `tfsdk:"external_id"`
	ExecutionRegions         types.List   `tfsdk:"execution_regions"`
	AWSAccountCustomerConfig types.Object `tfsdk:"aws_account_customer_config"`
	AWSOrgReadCustomerConfig types.Object `tfsdk:"aws_org_read_customer_config"`
}

// PlatformCustomerConfig is the model for customer config definitions.
type PlatformCustomerConfig struct {
	TerraformModule types.Object `tfsdk:"terraform_module"`
}

func (c PlatformCustomerConfig) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"terraform_module": basetypes.ObjectType{
			AttrTypes: TerraformModule{}.AttributeTypes(),
		},
	}
}

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
