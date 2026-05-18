// Copyright Stacklet, Inc. 2025, 2026

package datasources

import (
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type datasources struct {
	Released   []func() datasource.DataSource
	Unreleased []func() datasource.DataSource
}

// List returns available data source factories, optionally including unreleased ones.
func (d datasources) List(includeUnreleased bool) []func() datasource.DataSource {
	result := slices.Clone(d.Released)
	if includeUnreleased {
		result = append(result, d.Unreleased...)
	}
	return result
}

// Data sources registered with the provider.
var DataSources = datasources{
	Released: []func() datasource.DataSource{
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
		newFactory(&gcpIntegrationDataSource{}),
		newFactory(&gcpIntegrationSurfaceDataSource{}),
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
	},
}

func newFactory[T any, R interface {
	*T
	datasource.DataSource
}](_ R) func() datasource.DataSource {
	return func() datasource.DataSource {
		return R(new(T))
	}
}
