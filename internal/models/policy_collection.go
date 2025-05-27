// Copyright (c) 2025 - Stacklet, Inc.

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PolicyCollectionResource is the model for a policy collection resource.
type PolicyCollectionResource struct {
	ID             types.String `tfsdk:"id"`
	UUID           types.String `tfsdk:"uuid"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	CloudProvider  types.String `tfsdk:"cloud_provider"`
	AutoUpdate     types.Bool   `tfsdk:"auto_update"`
	System         types.Bool   `tfsdk:"system"`
	Dynamic        types.Bool   `tfsdk:"dynamic"`
	RepositoryUUID types.String `tfsdk:"repository_uuid"`
}

// PolicyCollectionDatasource is the model for a policy collection data source.
type PolicyCollectionDataSource PolicyCollectionResource
