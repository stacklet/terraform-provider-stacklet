// Copyright Stacklet, Inc. 2025, 2026

package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Resources returns all available resources.
func Resources() []func() resource.Resource {
	return []func() resource.Resource{
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
}
