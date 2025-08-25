// Copyright (c) 2025 - Stacklet, Inc.

package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var DATASOURCES = []func() datasource.DataSource{
	NewAccountDataSource,
	NewAccountGroupDataSource,
	NewBindingDataSource,
	NewConfigurationProfileAccountOwnersDataSource,
	NewConfigurationProfileEmailDataSource,
	NewConfigurationProfileJiraDataSource,
	NewConfigurationProfileResourceOwnerDataSource,
	NewConfigurationProfileServiceNowDataSource,
	NewConfigurationProfileSlackDataSource,
	NewConfigurationProfileSymphonyDataSource,
	NewConfigurationProfileTeamsDataSource,
	NewNotificationTemplateDataSource,
	NewPlatformDataSource,
	NewPolicyCollectionDataSource,
	NewPolicyDataSource,
	NewReportgroupDataSource,
	NewRepositoryDataSource,
}
