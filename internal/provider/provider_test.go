// Copyright Stacklet, Inc. 2025, 2026

package provider

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setHomeDir(t *testing.T) string {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	return homeDir
}

func makeStackletAdminDir(t *testing.T) string {
	homeDir := setHomeDir(t)
	stackletDir := path.Join(homeDir, ".stacklet")
	if err := os.MkdirAll(stackletDir, 0o755); err != nil {
		t.Fatal(err)
	}

	return stackletDir
}

func setupStackletAdminConfig(t *testing.T, endpoint string, apiKey string) {
	stackletDir := makeStackletAdminDir(t)

	if endpoint != "" {
		configFile := path.Join(stackletDir, "config.json")
		configData := map[string]string{"api": endpoint}
		configJSON, _ := json.Marshal(configData)
		if err := os.WriteFile(configFile, configJSON, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	if apiKey != "" {
		credsFile := path.Join(stackletDir, "credentials")
		if err := os.WriteFile(credsFile, []byte(apiKey), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetCredentials(t *testing.T) {
	envFull := providerEnv{
		Endpoint: "https://env-endpoint.example.com",
		APIKey:   "env-api-key",
	}
	configFull := providerModel{
		Endpoint: types.StringValue("https://config-endpoint.example.com"),
		APIKey:   types.StringValue("config-api-key"),
	}
	adminCLIFull := &credentials{
		Endpoint: "https://cli-endpoint.example.com",
		APIKey:   "cli-api-key",
	}

	tests := []struct {
		name       string
		env        providerEnv
		config     providerModel
		adminCLI   *credentials
		expected   credentials
		diagErrors []string
	}{
		{
			name:   "FromConfig",
			config: configFull,
			expected: credentials{
				Endpoint: "https://config-endpoint.example.com",
				APIKey:   "config-api-key",
			},
		},
		{
			name: "FromEnviron",
			env:  envFull,
			expected: credentials{
				Endpoint: "https://env-endpoint.example.com",
				APIKey:   "env-api-key",
			},
		},
		{
			name:     "FromAdminCLI",
			adminCLI: adminCLIFull,
			expected: credentials{
				Endpoint: "https://cli-endpoint.example.com",
				APIKey:   "cli-api-key",
			},
		},
		{
			name:   "ConfigOverridesEnviron",
			env:    envFull,
			config: configFull,
			expected: credentials{
				Endpoint: "https://config-endpoint.example.com",
				APIKey:   "config-api-key",
			},
		},
		{
			name:     "EnvironOverridesAdminCLI",
			env:      envFull,
			adminCLI: adminCLIFull,
			expected: credentials{
				Endpoint: "https://env-endpoint.example.com",
				APIKey:   "env-api-key",
			},
		},
		{
			name: "MixedSources",
			env: providerEnv{
				APIKey: "env-api-key",
			},
			adminCLI: &credentials{
				Endpoint: "https://cli-endpoint.example.com",
			},
			expected: credentials{
				Endpoint: "https://cli-endpoint.example.com",
				APIKey:   "env-api-key",
			},
		},
		{
			name: "IaCEndpointFromConfig",
			env: providerEnv{
				IaCEndpoint: "https://env-iac-endpoint.example.com",
			},
			config: providerModel{
				Endpoint:    types.StringValue("https://api.example.com"),
				IaCEndpoint: types.StringValue("https://config-iac-endpoint.example.com"),
				APIKey:      types.StringValue("config-api-key"),
			},
			expected: credentials{
				Endpoint:    "https://api.example.com",
				IaCEndpoint: "https://config-iac-endpoint.example.com",
				APIKey:      "config-api-key",
			},
		},
		{
			name: "IaCEndpointFromEnviron",
			env: providerEnv{
				IaCEndpoint: "https://env-iac-endpoint.example.com",
			},
			config: providerModel{
				Endpoint: types.StringValue("https://api.example.com"),
				APIKey:   types.StringValue("config-api-key"),
			},
			expected: credentials{
				Endpoint:    "https://api.example.com",
				IaCEndpoint: "https://env-iac-endpoint.example.com",
				APIKey:      "config-api-key",
			},
		},
		{
			name: "IaCEndpointDerived",
			adminCLI: &credentials{
				Endpoint: "https://api.example.com/",
				APIKey:   "cli-api-key",
			},
			expected: credentials{
				Endpoint:    "https://api.example.com/",
				IaCEndpoint: "https://api.sinistral.example.com/",
				APIKey:      "cli-api-key",
			},
		},
		{
			name: "UnknownIaCEndpoint",
			config: providerModel{
				Endpoint:    types.StringValue("https://api.example.com"),
				IaCEndpoint: types.StringUnknown(),
				APIKey:      types.StringValue("config-api-key"),
			},
			diagErrors: []string{"Unknown Stacklet IaC API Endpoint"},
		},
		{
			name: "MissingEndpoint",
			config: providerModel{
				APIKey: types.StringValue("config-api-key"),
			},
			diagErrors: []string{"Missing Stacklet API Endpoint"},
		},
		{
			name: "MissingAPIKey",
			config: providerModel{
				Endpoint: types.StringValue("https://config-endpoint.example.com"),
			},
			diagErrors: []string{"Missing Stacklet API key"},
		},
		{
			name: "MissingBoth",
			diagErrors: []string{
				"Missing Stacklet API Endpoint",
				"Missing Stacklet API key",
			},
		},
		{
			name: "UnknownEndpoint",
			config: providerModel{
				Endpoint: types.StringUnknown(),
				APIKey:   types.StringValue("config-api-key"),
			},
			diagErrors: []string{"Unknown Stacklet API Endpoint"},
		},
		{
			name: "UnknownAPIKey",
			config: providerModel{
				Endpoint: types.StringValue("https://config-endpoint.example.com"),
				APIKey:   types.StringUnknown(),
			},
			diagErrors: []string{"Unknown Stacklet API key"},
		},
		{
			name: "UnknownBoth",
			config: providerModel{
				Endpoint: types.StringUnknown(),
				APIKey:   types.StringUnknown(),
			},
			diagErrors: []string{
				"Unknown Stacklet API Endpoint",
				"Unknown Stacklet API key",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.adminCLI != nil {
				setupStackletAdminConfig(t, tc.adminCLI.Endpoint, tc.adminCLI.APIKey)
			} else {
				setHomeDir(t)
			}

			creds, diags := getCredentials(tc.config, tc.env)
			expectedErrors := len(tc.diagErrors)
			if expectedErrors == 0 {
				assert.Equal(t, tc.expected, creds)
			}
			require.Len(t, diags.Errors(), expectedErrors)
			for i, msg := range tc.diagErrors {
				assert.Contains(t, diags[i].Summary(), msg)
			}
		})
	}
}

func TestGetCredentials_InvalidAdminCLIConfigJSON(t *testing.T) {
	stackletDir := makeStackletAdminDir(t)
	configFile := path.Join(stackletDir, "config.json")
	if err := os.WriteFile(configFile, []byte("invalid json"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, diags := getCredentials(providerModel{}, providerEnv{})
	require.True(t, diags.HasError())
	assert.Contains(t, diags[0].Summary(), "Missing Stacklet API Endpoint")
}

func TestIaCEndpointFromAPIEndpoint(t *testing.T) {
	tests := []struct {
		endpoint string
		expected string
	}{
		{"https://api.example.com", "https://api.sinistral.example.com/"},
		{"https://api.example.com/", "https://api.sinistral.example.com/"},
		{"https://api.example.com/graphql", "https://api.sinistral.example.com/"},
		{"https://api.example.com:8443/", "https://api.sinistral.example.com:8443/"},
		{"https://example.com/", ""},
		{"https://myapi.example.com/", ""},
		{"", ""},
		{"://invalid", ""},
	}
	for _, tc := range tests {
		t.Run(tc.endpoint, func(t *testing.T) {
			assert.Equal(t, tc.expected, iacEndpointFromAPIEndpoint(tc.endpoint))
		})
	}
}
