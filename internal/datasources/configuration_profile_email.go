// Copyright (c) 2025 - Stacklet, Inc.

package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
)

var (
	_ datasource.DataSource = &configurationProfileEmailDataSource{}
)

func newConfigurationProfileEmailDataSource() datasource.DataSource {
	return &configurationProfileEmailDataSource{}
}

type configurationProfileEmailDataSource struct {
	apiDataSource
}

func (d *configurationProfileEmailDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configuration_profile_email"
}

func (d *configurationProfileEmailDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve information about the email configuration profile.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the configuration profile.",
				Computed:    true,
			},
			"profile": schema.StringAttribute{
				Description: "The profile name.",
				Computed:    true,
			},
			"from": schema.StringAttribute{
				Description: "The email from field value.",
				Computed:    true,
			},
			"ses_region": schema.StringAttribute{
				Description: "The SES region in use.",
				Computed:    true,
			},
			"smtp": schema.SingleNestedAttribute{
				Description: "SMTP configuration.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"server": schema.StringAttribute{
						Description: "SMTP server hostname or IP address.",
						Computed:    true,
					},
					"port": schema.StringAttribute{
						Description: "SMTP server port.",
						Computed:    true,
					},
					"ssl": schema.BoolAttribute{
						Description: "Whether SSL/TLS is enabled.",
						Computed:    true,
					},
					"username": schema.StringAttribute{
						Description: "Authentication username.",
						Computed:    true,
					},
					"password": schema.StringAttribute{
						Description: "Authentication password (encrypted).",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d *configurationProfileEmailDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.ConfigurationProfileEmailDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, err := d.api.ConfigurationProfile.ReadEmail(ctx)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(data.Update(ctx, *config)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
