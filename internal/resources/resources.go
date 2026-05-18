// Copyright Stacklet, Inc. 2025, 2026

package resources

import (
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type resources struct {
	Released   []func() resource.Resource
	Unreleased []func() resource.Resource
}

// List returns available resource factories, optionally including unreleased ones.
func (r resources) List(includeUnreleased bool) []func() resource.Resource {
	result := slices.Clone(r.Released)
	if includeUnreleased {
		result = append(result, r.Unreleased...)
	}
	return result
}

// Resources registered with the provider.
var Resources = resources{
	Released: []func() resource.Resource{
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
		newFactory(&gcpIntegrationResource{}),
		newFactory(&notificationTemplateResource{}),
		newFactory(&policyCollectionMappingResource{}),
		newFactory(&policyCollectionResource{}),
		newFactory(&reportGroupResource{}),
		newFactory(&repositoryResource{}),
		newFactory(&roleAssignmentResource{}),
		newFactory(&ssoGroupResource{}),
		newFactory(&userResource{}),
	},
}

func newFactory[T any, R interface {
	*T
	resource.Resource
}](_ R) func() resource.Resource {
	return func() resource.Resource {
		return R(new(T))
	}
}
