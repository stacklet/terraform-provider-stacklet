// Copyright (c) 2025 - Stacklet, Inc.

package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var RESOURCES = []func() resource.Resource{
	NewAccountDiscoveryAWSResource,
	NewAccountDiscoveryAzureResource,
	NewAccountDiscoveryGCPResource,
	NewAccountGroupMappingResource,
	NewAccountGroupResource,
	NewAccountResource,
	NewBindingResource,
	NewConfigurationProfileAccountOwnersResource,
	NewConfigurationProfileEmailResource,
	NewConfigurationProfileJiraResource,
	NewConfigurationProfileResourceOwnerResource,
	NewConfigurationProfileServiceNowResource,
	NewConfigurationProfileSlackResource,
	NewConfigurationProfileSymphonyResource,
	NewConfigurationProfileTeamsResource,
	NewNotificationTemplateResource,
	NewPolicyCollectionMappingResource,
	NewPolicyCollectionResource,
	NewReportGroupResource,
	NewRepositoryResource,
}
