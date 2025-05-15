package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var RESOURCES = []func() resource.Resource{
	NewAccountResource,
	NewAccountGroupResource,
	NewAccountGroupMappingResource,
	NewPolicyCollectionItemResource,
	NewPolicyCollectionResource,
	NewBindingResource,
	NewAccountDiscoveryResource,
	NewRepositoryResource,
}
