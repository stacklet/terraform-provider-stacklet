// Copyright Stacklet, Inc. 2025, 2026

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
)

var _ datasource.DataSource = &iacProjectDataSource{}

type iacProjectDataSource struct {
	iacAPIDataSource
}

func (d *iacProjectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iac_project"
}

func (d *iacProjectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve information about an IaC project by name.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the project.",
				Required:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "The client ID of the project credentials, if any.",
				Computed:    true,
			},
			"collections": schema.SetAttribute{
				ElementType: types.StringType,
				Description: "The names of the policy collections used by the project.",
				Computed:    true,
			},
			"read_groups": schema.SetAttribute{
				ElementType: types.StringType,
				Description: "The names of the groups with read access to the project.",
				Computed:    true,
			},
		},
	}
}

func (d *iacProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.IaCProjectDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := d.api.Project.Read(ctx, data.Name.ValueString())
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(data.Update(ctx, project)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
