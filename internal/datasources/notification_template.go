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
	_ datasource.DataSource = &notificationTemplateDataSource{}
)

func newNotificationTemplateDataSource() datasource.DataSource {
	return &notificationTemplateDataSource{}
}

type notificationTemplateDataSource struct {
	apiDataSource
}

func (d *notificationTemplateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_template"
}

func (d *notificationTemplateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetch a notification template by name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The GraphQL Node ID of the template.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the template.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the template.",
				Optional:    true,
				Computed:    true,
			},
			"transport": schema.StringAttribute{
				Description: "The notification transport the template is for.",
				Optional:    true,
				Computed:    true,
			},
			"content": schema.StringAttribute{
				Description: "The template content.",
				Computed:    true,
			},
		},
	}
}

func (d *notificationTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.NotificationTemplateDataSource
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	template, err := d.api.Template.Read(ctx, data.Name.ValueString())
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(data.Update(template)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
