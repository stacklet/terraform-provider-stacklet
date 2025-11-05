// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/typehelpers"
)

// PlatformDataSource is the model for the platform data source.
type PlatformDataSource struct {
	ID                       types.String `tfsdk:"id"`
	ExternalID               types.String `tfsdk:"external_id"`
	ExecutionRegions         types.List   `tfsdk:"execution_regions"`
	AWSAccountCustomerConfig types.Object `tfsdk:"aws_account_customer_config"`
	AWSOrgReadCustomerConfig types.Object `tfsdk:"aws_org_read_customer_config"`
}

func (m *PlatformDataSource) Update(ctx context.Context, platform *api.Platform) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(platform.ID)
	m.ExternalID = types.StringPointerValue(platform.ExternalID)
	m.ExecutionRegions = typehelpers.StringsList(platform.ExecutionRegions)

	awsAccountCustomerConfig, d := m.getCustomerConfig(ctx, platform.AWSAccountCustomerConfig)
	diags.Append(d...)
	m.AWSAccountCustomerConfig = awsAccountCustomerConfig

	awsOrgReadCustomerConfig, d := m.getCustomerConfig(ctx, platform.AWSOrgReadCustomerConfig)
	diags.Append(d...)
	m.AWSOrgReadCustomerConfig = awsOrgReadCustomerConfig

	return diags
}

func (m PlatformDataSource) getCustomerConfig(ctx context.Context, config api.PlatformCustomerConfig) (types.Object, diag.Diagnostics) {
	terraformModule, diags := typehelpers.ObjectValue(
		ctx,
		&config.TerraformModule,
		func() (*TerraformModule, diag.Diagnostics) {
			return &TerraformModule{
				RepositoryURL: types.StringValue(config.TerraformModule.RepositoryURL),
				Source:        types.StringValue(config.TerraformModule.Source),
				VariablesJSON: types.StringValue(config.TerraformModule.VariablesJSON),
				Version:       types.StringPointerValue(config.TerraformModule.Version),
			}, nil
		},
	)
	if diags.HasError() {
		return types.ObjectNull(PlatformCustomerConfig{}.AttributeTypes()), diags
	}

	return typehelpers.ObjectValue(
		ctx,
		&config,
		func() (*PlatformCustomerConfig, diag.Diagnostics) {
			return &PlatformCustomerConfig{
				TerraformModule: terraformModule,
			}, nil
		},
	)
}

// PlatformCustomerConfig is the model for customer config definitions.
type PlatformCustomerConfig struct {
	TerraformModule types.Object `tfsdk:"terraform_module"`
}

func (c PlatformCustomerConfig) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"terraform_module": types.ObjectType{
			AttrTypes: TerraformModule{}.AttributeTypes(),
		},
	}
}
