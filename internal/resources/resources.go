package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var RESOURCES = []func() resource.Resource{
	NewAccountDiscoveryAzureResource,
	NewAccountDiscoveryGCPResource,
	NewAccountGroupMappingResource,
	NewAccountGroupResource,
	NewAccountResource,
	NewBindingResource,
	NewPolicyCollectionMappingResource,
	NewPolicyCollectionResource,
	NewRepositoryResource,
}
