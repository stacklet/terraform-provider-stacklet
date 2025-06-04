// Copyright (c) 2025 - Stacklet, Inc.

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
	"github.com/stacklet/terraform-provider-stacklet/internal/providerdata"
	tftypes "github.com/stacklet/terraform-provider-stacklet/internal/types"
)

var (
	_ datasource.DataSource = &bindingExecutionConfigDataSource{}
)

func NewBindingExecutionConfigDataSource() datasource.DataSource {
	return &bindingExecutionConfigDataSource{}
}

type bindingExecutionConfigDataSource struct {
	api *api.API
}

func (d *bindingExecutionConfigDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_binding_execution_config"
}

func (d *bindingExecutionConfigDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetch execution configuration for a binding.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL ID of the binding the configuration is for.",
				Computed:    true,
			},
			"binding_uuid": schema.StringAttribute{
				Description: "The UUID of the binding the configuration is for.",
				Required:    true,
			},
			"dry_run": schema.BoolAttribute{
				Description: "Whether the binding is run in with action disabled (in information mode).",
				Optional:    true,
				Computed:    true,
			},
			"variables": schema.StringAttribute{
				Description: "JSON-encoded dictionary of values used for policy templating.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (d *bindingExecutionConfigDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if pd, err := providerdata.GetDataSourceProviderData(req); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
	} else if pd != nil {
		d.api = pd.API
	}
}

func (d *bindingExecutionConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.BindingExecutionConfigDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	executionConfig, err := d.api.BindingExecutionConfig.Read(ctx, data.BindingUUID.ValueString())
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	data.ID = data.BindingUUID // since it's not a separate GraphQL entity, use the binding UUID as ID too
	data.DryRun = types.BoolValue(executionConfig.DryRunDefault())
	variablesString, err := tftypes.JSONString(executionConfig.Variables)
	if err != nil {
		resp.Diagnostics.AddError("Invalid content for variables", err.Error())
		return
	}
	data.Variables = variablesString
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
