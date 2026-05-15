// Copyright Stacklet, Inc. 2025, 2026

package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// DataSources returns all available datasources.
func DataSources(includeUnreleased bool) []func() datasource.DataSource {
	dataSources := []func() datasource.DataSource{
		newFactory(&accountDataSource{}),
		newFactory(&accountGroupDataSource{}),
		newFactory(&bindingDataSource{}),
		newFactory(&configurationProfileAccountOwnersDataSource{}),
		newFactory(&configurationProfileEmailDataSource{}),
		newFactory(&configurationProfileJiraDataSource{}),
		newFactory(&configurationProfileMSTeamsDataSource{}),
		newFactory(&configurationProfileResourceOwnerDataSource{}),
		newFactory(&configurationProfileServiceNowDataSource{}),
		newFactory(&configurationProfileSlackDataSource{}),
		newFactory(&configurationProfileSymphonyDataSource{}),
		newFactory(&msteamsIntegrationSurfaceDataSource{}),
		newFactory(&notificationTemplateDataSource{}),
		newFactory(&platformDataSource{}),
		newFactory(&policyCollectionDataSource{}),
		newFactory(&policyDataSource{}),
		newFactory(&reportGroupDataSource{}),
		newFactory(&repositoryDataSource{}),
		newFactory(&roleAssignmentsDataSource{}),
		newFactory(&roleDataSource{}),
		newFactory(&ssoGroupDataSource{}),
		newFactory(&userDataSource{}),
	}
	if includeUnreleased {
		dataSources = append(
			dataSources,
			newFactory(&gcpIntegrationDataSource{}),
			newFactory(&gcpIntegrationSurfaceDataSource{}),
		)
	}
	return dataSources
}

func newFactory[T any, R interface {
	*T
	datasource.DataSource
}](_ R) func() datasource.DataSource {
	return func() datasource.DataSource {
		return R(new(T))
	}
}
