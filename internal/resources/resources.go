package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var RESOURCES = []func() resource.Resource{
	NewAccountResource,
	NewAccountGroupResource,
	NewAccountGroupItemResource,
	NewPolicyCollectionItemResource,
	NewPolicyCollectionResource,
	NewBindingResource,
	NewAccountDiscoveryResource,
	NewSSOGroupResource,
	NewRepositoryResource,
}
