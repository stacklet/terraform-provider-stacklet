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

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &iacPolicySourceDataSource{}

// iacPolicySourceDataSource defines the data source implementation.
type iacPolicySourceDataSource struct {
	iacAPIDataSource
}

func (d *iacPolicySourceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iac_policy_source"
}

func (d *iacPolicySourceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetch a Stacklet IaC policy source, a repository containing IaC policies.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the policy source.",
				Required:    true,
			},
			"url": schema.StringAttribute{
				Description: "The URL of the remote repository.",
				Computed:    true,
			},
			"branch": schema.StringAttribute{
				Description: "The branch scanned for policies. If null, the default branch is used.",
				Computed:    true,
			},
			"paths": schema.SetAttribute{
				Description: "The directories scanned for policies. If null, the whole repository is scanned.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"extensions": schema.SetAttribute{
				Description: "The extensions of policy files.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"auto_scan": schema.BoolAttribute{
				Description: "Whether the repository is periodically scanned for policy changes.",
				Computed:    true,
			},
			"webhook_url": schema.StringAttribute{
				Description: "The URL of the webhook which triggers repository scans.",
				Computed:    true,
			},
			"token_auth": schema.SingleNestedAttribute{
				Description: "Authentication to the remote repository with a user and password or token, if configured.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"user": schema.StringAttribute{
						Description: "The user with access to the remote repository.",
						Computed:    true,
					},
				},
			},
			"ssh_auth": schema.SingleNestedAttribute{
				Description: "Authentication to the remote repository with an SSH key, if configured.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"user": schema.StringAttribute{
						Description: "The user with access to the remote repository.",
						Computed:    true,
					},
					"has_passphrase": schema.BoolAttribute{
						Description: "Whether the SSH private key has a passphrase.",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d *iacPolicySourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.IaCPolicySourceDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	source, err := d.api.PolicySource.Read(ctx, data.Name.ValueString())
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(data.Update(ctx, source)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
