// Copyright Stacklet, Inc. 2025, 2026

package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

// GCPIntegrationSurfaceTrustAWS holds details for the AWS trust configuration for GCP integration.
type GCPIntegrationSurfaceTrustAWS struct {
	AccountID         types.String `tfsdk:"account_id"`
	AssetDBRoleName   types.String `tfsdk:"assetdb_role_name"`
	CostQueryRoleName types.String `tfsdk:"cost_query_role_name"`
	ExecutionRoleName types.String `tfsdk:"execution_role_name"`
	PlatformRoleName  types.String `tfsdk:"platform_role_name"`
}

// GCPIntegrationSurfaceAWSRelay holds details for the AWS relay configuration for GCP integration.
type GCPIntegrationSurfaceAWSRelay struct {
	BusARN  types.String `tfsdk:"bus_arn"`
	RoleARN types.String `tfsdk:"role_arn"`
}

// GCPIntegrationSurfaceDataSource is the model for the GCP integration surface data source.
type GCPIntegrationSurfaceDataSource struct {
	TrustAWS *GCPIntegrationSurfaceTrustAWS `tfsdk:"trust_aws"`
	AWSRelay *GCPIntegrationSurfaceAWSRelay `tfsdk:"aws_relay"`
}

func (m *GCPIntegrationSurfaceDataSource) Update(_ context.Context, surface *api.GCPIntegrationSurface) diag.Diagnostics {
	var diags diag.Diagnostics

	m.TrustAWS = &GCPIntegrationSurfaceTrustAWS{
		AccountID:         types.StringValue(surface.TrustAws.AccountID),
		AssetDBRoleName:   types.StringValue(surface.TrustAws.AssetdbRoleName),
		CostQueryRoleName: types.StringValue(surface.TrustAws.CostQueryRoleName),
		ExecutionRoleName: types.StringValue(surface.TrustAws.ExecutionRoleName),
		PlatformRoleName:  types.StringValue(surface.TrustAws.PlatformRoleName),
	}
	m.AWSRelay = &GCPIntegrationSurfaceAWSRelay{
		BusARN:  types.StringValue(surface.AwsRelay.BusArn),
		RoleARN: types.StringValue(surface.AwsRelay.RoleArn),
	}

	return diags
}
