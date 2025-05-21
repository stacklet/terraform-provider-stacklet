package schemavalidate

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/stacklet/terraform-provider-stacklet/internal/api"
)

// OneOfCloudProviders is a validator that checks that value is one of the supported cloud providers.
func OneOfCloudProviders() validator.String {
	allProviders := make([]string, len(api.CLOUD_PROVIDERS))
	for i, provider := range api.CLOUD_PROVIDERS {
		allProviders[i] = string(provider)
	}
	return stringvalidator.OneOf(allProviders...)
}
