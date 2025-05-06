package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

// CloudProvider represents a cloud service provider in Stacklet
type CloudProvider string

const (
	CloudProviderAWS        CloudProvider = "AWS"
	CloudProviderAzure      CloudProvider = "AZURE"
	CloudProviderGCP        CloudProvider = "GCP"
	CloudProviderKubernetes CloudProvider = "KUBERNETES"
	CloudProviderTencent    CloudProvider = "TENCENT"
)

// NewCloudProvider creates a CloudProvider from a string and validates it.
func NewCloudProvider(s string) (CloudProvider, error) {
	p := CloudProvider(strings.ToUpper(s))
	return p, p.Validate()
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (cp *CloudProvider) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*cp = CloudProvider(strings.ToUpper(s))
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (cp CloudProvider) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(cp))
}

// String implements the fmt.Stringer interface
func (cp CloudProvider) String() string {
	return string(cp)
}

// Validate returns an error if the cloud provider is not valid
func (cp CloudProvider) Validate() error {
	switch cp {
	case CloudProviderAWS, CloudProviderAzure, CloudProviderGCP, CloudProviderKubernetes, CloudProviderTencent:
		return nil
	default:
		return fmt.Errorf("invalid cloud provider: %s", cp)
	}
}
