package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var DATASOURCES = []func() datasource.DataSource{
	NewAccountDataSource,
	NewAccountGroupDataSource,
	NewPolicyDataSource,
	NewPolicyCollectionDataSource,
	NewBindingDataSource,
	NewAccountDiscoveryDataSource,
	NewSSOGroupDataSource,
	NewRepositoryDataSource,
}
