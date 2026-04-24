// Copyright Stacklet, Inc. 2025, 2026



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
		newMSTeamsIntegrationSurfaceDataSource,
		newNotificationTemplateDataSource,
		newPlatformDataSource,
		newPolicyCollectionDataSource,
		newPolicyDataSource,
		newReportgroupDataSource,
		newRepositoryDataSource,
		newRoleDataSource,
		newRoleAssignmentsDataSource,
		newSSOGroupDataSource,
		newUserDataSource,
	}
}
