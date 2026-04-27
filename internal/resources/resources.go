// Copyright Stacklet, Inc. 2025, 2026

package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Resources returns all available resources.
func Resources(includeUnreleased bool) []func() resource.Resource {
	resources := []func() resource.Resource{
		newAccountDiscoveryAWSResource,
		newAccountDiscoveryAzureResource,
		newAccountDiscoveryGCPResource,
		newAccountGroupMappingResource,
		newAccountGroupResource,
		newAccountResource,
		newBindingResource,
		newConfigurationProfileAccountOwnersResource,
		newConfigurationProfileEmailResource,
		newConfigurationProfileJiraResource,
		newConfigurationProfileMSTeamsResource,
		newConfigurationProfileResourceOwnerResource,
		newConfigurationProfileServiceNowResource,
		newConfigurationProfileSlackResource,
		newConfigurationProfileSymphonyResource,
		newNotificationTemplateResource,
		newPolicyCollectionMappingResource,
		newPolicyCollectionResource,
		newReportGroupResource,
		newRepositoryResource,
		newRoleAssignmentResource,
		newSSOGroupResource,
		newUserResource,
	}
	unreleasedResources := []func() resource.Resource{}
	if includeUnreleased {
		resources = append(resources, unreleasedResources...)
	}
	return resources
}
