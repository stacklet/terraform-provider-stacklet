// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

// NotificationTemplateResource is the model for a notification template resource.
type NotificationTemplateResource struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Transport   types.String `tfsdk:"transport"`
	Content     types.String `tfsdk:"content"`
}

func (m *NotificationTemplateResource) Update(template *api.Template) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.StringValue(template.ID)
	m.Name = types.StringValue(template.Name)
	m.Description = types.StringPointerValue(template.Description)
	m.Transport = types.StringPointerValue(template.Transport)
	m.Content = types.StringValue(template.Content)

	return diags
}

// NotificationTemplateDataSource is the model for a notification template data source.
type NotificationTemplateDataSource struct {
	NotificationTemplateResource
}
