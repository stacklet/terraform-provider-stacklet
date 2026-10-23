// Copyright Stacklet, Inc. 2025, 2026

package providerdata

import (
	"fmt"

	"github.com/stacklet/terraform-provider-stacklet/internal/api/iac"
	"github.com/stacklet/terraform-provider-stacklet/internal/api/platform"
)

// providerData holds shared data for available in requests.
type providerData struct {
	PlatformAPI *platform.API
	IaCAPI      *iac.API
}

// New returns configured provider data.
func New(config platform.ClientConfig, iacConfig iac.ClientConfig) *providerData {
	return &providerData{
		PlatformAPI: platform.New(config),
		IaCAPI:      iac.New(iacConfig),
	}
}

// APIClient is the set of API clients available from provider data.
type APIClient interface {
	*platform.API | *iac.API
}

type providerDataError struct {
	Kind         string
	ProviderData any
}

func (e providerDataError) Summary() string {
	return fmt.Sprintf("Unexpected %s configure type", e.Kind)
}

func (e providerDataError) Error() string {
	return fmt.Sprintf("Expected *ProviderData, got: %T. Please report this issue to the provider developers.", e.ProviderData)
}

// GetAPI returns the API client of the requested type from the provider data
// of a configure request, or nil if provider data is not set.
//
// kind describes the type of object being configured, for error reporting.
func GetAPI[T APIClient](kind string, data any) (T, error) {
	var api T
	if data == nil {
		return api, nil
	}
	pd, ok := data.(*providerData)
	if !ok {
		return api, providerDataError{
			Kind:         kind,
			ProviderData: data,
		}
	}
	switch p := any(&api).(type) {
	case **platform.API:
		*p = pd.PlatformAPI
	case **iac.API:
		*p = pd.IaCAPI
	}
	return api, nil
}
