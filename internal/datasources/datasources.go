// Copyright (c) 2025 - Stacklet, Inc.

package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// DataSources returns all available datasources..
func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newAccountDataSource,
		newAccountGroupDataSource,
		newBindingDataSource,
		newConfigurationProfileAccountOwnersDataSource,
		newConfigurationProfileEmailDataSource,
		newConfigurationProfileJiraDataSource,
		newConfigurationProfileMSTeamsDataSource,
		newConfigurationProfileResourceOwnerDataSource,
		newConfigurationProfileServiceNowDataSource,
		newConfigurationProfileSlackDataSource,
		newConfigurationProfileSymphonyDataSource,
		newNotificationTemplateDataSource,
		newPlatformDataSource,
		newPolicyCollectionDataSource,
		newPolicyDataSource,
		newReportgroupDataSource,
		newRepositoryDataSource,
	}
}
