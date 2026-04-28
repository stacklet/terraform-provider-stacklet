// Copyright Stacklet, Inc. 2025, 2026

package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// DataSources returns all available datasources..
func DataSources(includeUnreleased bool) []func() datasource.DataSource {
	dataSources := []func() datasource.DataSource{
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
		newMSTeamsIntegrationSurfaceDataSource,
		newNotificationTemplateDataSource,
		newPlatformDataSource,
		newPolicyCollectionDataSource,
		newPolicyDataSource,
		newReportgroupDataSource,
		newRepositoryDataSource,
		newRoleAssignmentsDataSource,
		newRoleDataSource,
		newSSOGroupDataSource,
		newUserDataSource,
	}
	unreleasedDataSources := []func() datasource.DataSource{
		newGCPIntegrationDataSource,
		newGCPIntegrationSurfaceDataSource,
	}
	if includeUnreleased {
		dataSources = append(dataSources, unreleasedDataSources...)
	}
	return dataSources
}
