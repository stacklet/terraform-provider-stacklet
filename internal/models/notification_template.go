// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NotificationTemplateResource is the model for a notification template resource.
type NotificationTemplateResource struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Transport   types.String `tfsdk:"transport"`
	Content     types.String `tfsdk:"content"`
}

// NotificationTemplateDataSource is the model for a notification template data source.
type NotificationTemplateDataSource NotificationTemplateResource
