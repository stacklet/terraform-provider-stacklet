package api

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	CloudProviderAWS          = CloudProvider("AWS")
	CloudProviderAzure        = CloudProvider("Azure")
	CloudProviderGCP          = CloudProvider("GCP")
	CloudProviderKubernetes   = CloudProvider("Kubernetes")
	CloudProviderTencentCloud = CloudProvider("TencentCloud")
)

// CLOUD_PROVIDERS is the list of all supported cloud providers.
var CLOUD_PROVIDERS = []CloudProvider{
	CloudProviderAWS,
	CloudProviderAzure,
	CloudProviderGCP,
	CloudProviderKubernetes,
	CloudProviderTencentCloud,
}

// CloudProvider represents a cloud service provider in Stacklet.
type CloudProvider string

// UnmarshalJSON implements the json.Unmarshaler interface.
func (cp *CloudProvider) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*cp = CloudProvider(s)
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (cp CloudProvider) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(cp))
}

// String implements the fmt.Stringer interface.
func (cp CloudProvider) String() string {
	return string(cp)
}

// NullableString converts a types.String to a string pointer which can be null.
func NullableString(s types.String) *string {
	if s.IsNull() {
		return nil
	}

	str := s.ValueString()
	return &str
}

// NullableBool converts a types.Bool to a bool pointer which can be null.
func NullableBool(b types.Bool) *bool {
	if b.IsNull() {
		return nil
	}

	bv := b.ValueBool()
	return &bv
}

// StringsList concerts a types.List to a list of strings.
func StringsList(l types.List) []string {
	if l.IsNull() || l.IsUnknown() {
		return nil
	}
	elements := l.Elements()
	sl := make([]string, len(elements))
	for i, element := range elements {
		str, _ := element.(types.String)
		sl[i] = str.ValueString()
	}
	return sl
}
