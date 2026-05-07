// Copyright Stacklet, Inc. 2025, 2026

package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Resources returns all available resources.
func Resources(includeUnreleased bool) []func() resource.Resource {
	resources := []func() resource.Resource{
		newFactory(&accountDiscoveryAWSResource{}),
		newFactory(&accountDiscoveryAzureResource{}),
		newFactory(&accountDiscoveryGCPResource{}),
		newFactory(&accountGroupMappingResource{}),
		newFactory(&accountGroupResource{}),
		newFactory(&accountResource{}),
		newFactory(&bindingResource{}),
		newFactory(&configurationProfileAccountOwnersResource{}),
		newFactory(&configurationProfileEmailResource{}),
		newFactory(&configurationProfileJiraResource{}),
		newFactory(&configurationProfileMSTeamsResource{}),
		newFactory(&configurationProfileResourceOwnerResource{}),
		newFactory(&configurationProfileServiceNowResource{}),
		newFactory(&configurationProfileSlackResource{}),
		newFactory(&configurationProfileSymphonyResource{}),
		newFactory(&notificationTemplateResource{}),
		newFactory(&policyCollectionMappingResource{}),
		newFactory(&policyCollectionResource{}),
		newFactory(&reportGroupResource{}),
		newFactory(&repositoryResource{}),
		newFactory(&roleAssignmentResource{}),
		newFactory(&ssoGroupResource{}),
		newFactory(&userResource{}),
	}
	if includeUnreleased {
		resources = append(resources, newFactory(&gcpIntegrationResource{}))
	}
	return resources
}

func newFactory[T any, R interface {
	*T
	resource.Resource
}](_ R) func() resource.Resource {
	return func() resource.Resource {
		return R(new(T))
	}
}
