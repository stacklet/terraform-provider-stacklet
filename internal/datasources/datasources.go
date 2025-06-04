// Copyright (c) 2025 - Stacklet, Inc.

package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var DATASOURCES = []func() datasource.DataSource{
	NewAccountDataSource,
	NewAccountGroupDataSource,
	NewBindingDataSource,
	NewBindingExecutionConfigDataSource,
	NewPolicyCollectionDataSource,
	NewPolicyDataSource,
	NewRepositoryDataSource,
}
