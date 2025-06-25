// Copyright (c) 2025 - Stacklet, Inc.

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
)

var (
	_ datasource.DataSource = &platformDataSource{}
)

func NewPlatformDataSource() datasource.DataSource {
	return &platformDataSource{}
}

type platformDataSource struct {
	api *api.API
}

func (d *platformDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_platform"
}

func (d *platformDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve information about the Stacklet platform.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID.",
				Computed:    true,
			},
			"external_id": schema.StringAttribute{
				Description: "The external ID for the deployment.",
				Computed:    true,
			},
			"execution_regions": schema.ListAttribute{
				Description: "List of regions for which execution is enabled.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"aws_account_customer_config": schema.SingleNestedAttribute{
				Description: "Customer configuration for AWS accounts.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"terraform_module": schema.SingleNestedAttribute{
						Description: "Terraform module configuration for account setup.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"repository_url": schema.StringAttribute{
								Description: "Module repository URL.",
								Computed:    true,
							},
							"source": schema.StringAttribute{
								Description: "Module source.",
								Computed:    true,
							},
							"variables_json": schema.StringAttribute{
								Description: "JSON-encoded variables for module configuration.",
								Computed:    true,
							},
						},
					},
				},
			},
			"aws_org_read_customer_config": schema.SingleNestedAttribute{
				Description: "Customer configuration for AWS organization read access.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"terraform_module": schema.SingleNestedAttribute{
						Description: "Terraform module configuration for organization read access setup.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"repository_url": schema.StringAttribute{
								Description: "Module repository URL.",
								Computed:    true,
							},
							"source": schema.StringAttribute{
								Description: "Module source.",
								Computed:    true,
							},
							"variables_json": schema.StringAttribute{
								Description: "JSON-encoded variables for module configuration.",
								Computed:    true,
							},
						},
					},
				},
			},
		},
	}
}

func (d *platformDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if pd, err := providerdata.GetDataSourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		d.api = pd.API
	}
}

func (d *platformDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.PlatformDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	platform, err := d.api.Platform.Read(ctx)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	data.ID = types.StringValue(platform.ID)
	data.ExternalID = types.StringPointerValue(platform.ExternalID)
	data.ExecutionRegions = tftypes.StringsList(platform.ExecutionRegions)
	awsAccountCustomerConfig, diags := d.getCustomerConfig(ctx, platform.AWSAccountCustomerConfig)
	resp.Diagnostics.Append(diags...)
	data.AWSAccountCustomerConfig = awsAccountCustomerConfig
	awsOrgReadCustomerConfig, diags := d.getCustomerConfig(ctx, platform.AWSOrgReadCustomerConfig)
	resp.Diagnostics.Append(diags...)
	data.AWSOrgReadCustomerConfig = awsOrgReadCustomerConfig

	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d platformDataSource) getCustomerConfig(ctx context.Context, config api.PlatformCustomerConfig) (basetypes.ObjectValue, diag.Diagnostics) {
	terraformModule, diags := tftypes.ObjectValue(
		ctx,
		&config.TerraformModule,
		func() (*models.TerraformModule, diag.Diagnostics) {
			return &models.TerraformModule{
				RepositoryURL: types.StringValue(config.TerraformModule.RepositoryURL),
				Source:        types.StringValue(config.TerraformModule.Source),
				VariablesJSON: types.StringValue(config.TerraformModule.VariablesJSON),
			}, nil
		},
	)
	if diags.HasError() {
		return basetypes.NewObjectNull(models.PlatformCustomerConfig{}.AttributeTypes()), diags
	}

	return tftypes.ObjectValue(
		ctx,
		&config,
		func() (*models.PlatformCustomerConfig, diag.Diagnostics) {
			return &models.PlatformCustomerConfig{
				TerraformModule: terraformModule,
			}, nil
		},
	)
}
