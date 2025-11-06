// Copyright (c) 2025 - Stacklet, Inc.

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
	}
}
