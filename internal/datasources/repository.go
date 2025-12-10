// Copyright (c) 2025 - Stacklet, Inc.

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &repositoryDataSource{}

func newRepositoryDataSource() datasource.DataSource {
	return &repositoryDataSource{}
}

// repositoryDataSource defines the data source implementation.
type repositoryDataSource struct {
	apiDataSource
}

func (d *repositoryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository"
}

func (d *repositoryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetch a Stacklet repository.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL ID of the repository.",
				Computed:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "The UUID of the repository.",
				Optional:    true,
				Computed:    true,
			},
			"url": schema.StringAttribute{
				Description: "The URL of the repository.",
				Optional:    true,
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the repository.",
				Computed:    true,
			},
			"webhook_url": schema.StringAttribute{
				Description: "The URL of the webhook which triggers repository scans.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description of the repository.",
				Computed:    true,
			},
			"auth_user": schema.StringAttribute{
				Description: "The user used to access the repository.",
				Computed:    true,
			},
			"has_auth_token": schema.BoolAttribute{
				Description: "Whether the repository has an auth token configured.",
				Computed:    true,
			},
			"ssh_public_key": schema.StringAttribute{
				Description: "If has_ssh_private_key, identifies that SSH private key.",
				Computed:    true,
			},
			"has_ssh_private_key": schema.BoolAttribute{
				Description: "Whether the repository has an SSH private key configured.",
				Computed:    true,
			},
			"has_ssh_passphrase": schema.BoolAttribute{
				Description: "Whether the repository has an SSH passphrase configured.",
				Computed:    true,
			},
			"system": schema.BoolAttribute{
				Description: "Whether this is a system repository (not user editable).",
				Computed:    true,
			},
			"role_assignment_target": schema.StringAttribute{
				Description: "The target identifier for role assignments (e.g., 'repository:uuid'). Use this value when assigning roles to this repository.",
				Computed:    true,
			},
		},
	}
}

func (d *repositoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.RepositoryDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Find by UUID, looking that up by URL if necessary.
	var uuid string
	if data.UUID.IsNull() || data.UUID.IsUnknown() {
		var err error
		uuid, err = d.api.Repository.FindByURL(ctx, data.URL.ValueString())
		if err != nil {
			errors.AddDiagError(&resp.Diagnostics, err)
			return
		}
	} else {
		uuid = data.UUID.ValueString()
	}
	repo, err := d.api.Repository.Read(ctx, uuid)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(data.Update(repo)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
