package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var RESOURCES = []func() resource.Resource{
	NewAccountResource,
	NewAccountGroupResource,
	NewAccountGroupMappingResource,
	NewPolicyCollectionMappingResource,
	NewPolicyCollectionResource,
	NewBindingResource,
	NewAccountDiscoveryResource,
	NewRepositoryResource,
}
