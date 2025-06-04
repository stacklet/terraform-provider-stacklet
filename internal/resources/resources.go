// Copyright (c) 2025 - Stacklet, Inc.

package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var RESOURCES = []func() resource.Resource{
	NewAccountDiscoveryAzureResource,
	NewAccountDiscoveryAWSResource,
	NewAccountDiscoveryGCPResource,
	NewAccountGroupMappingResource,
	NewAccountGroupResource,
	NewAccountResource,
	NewBindingResource,
	NewBindingExecutionConfigResource,
	NewPolicyCollectionMappingResource,
	NewPolicyCollectionResource,
	NewRepositoryResource,
}
