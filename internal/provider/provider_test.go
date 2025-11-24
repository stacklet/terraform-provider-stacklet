// Copyright (c) 2025 - Stacklet, Inc.

package provider

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func setHomeDir(t *testing.T) string {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	return homeDir
}

func makeStackletAdminDir(t *testing.T) string {
	homeDir := setHomeDir(t)
	stackletDir := path.Join(homeDir, ".stacklet")
	if err := os.MkdirAll(stackletDir, 0755); err != nil {
		t.Fatal(err)
	}

	return stackletDir
}

func setupStackletAdminConfig(t *testing.T, endpoint, apiKey string) {
	stackletDir := makeStackletAdminDir(t)

	if endpoint != "" {
		configFile := path.Join(stackletDir, "config.json")
		configData := map[string]string{"api": endpoint}
		configJSON, _ := json.Marshal(configData)
		if err := os.WriteFile(configFile, configJSON, 0644); err != nil {
			t.Fatal(err)
		}
	}

	if apiKey != "" {
		credsFile := path.Join(stackletDir, "credentials")
		if err := os.WriteFile(credsFile, []byte(apiKey), 0644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestGetCredentials_FromConfig(t *testing.T) {
	config := stackletProviderModel{
		Endpoint: types.StringValue("https://config-endpoint.example.com"),
		APIKey:   types.StringValue("config-api-key"),
	}

	creds := getCredentials(config)
	assert.Equal(t, "https://config-endpoint.example.com", creds.Endpoint)
	assert.Equal(t, "config-api-key", creds.APIKey)
}

func TestGetCredentials_FromEnvVars(t *testing.T) {
	t.Setenv("STACKLET_ENDPOINT", "https://env-endpoint.example.com")
	t.Setenv("STACKLET_API_KEY", "env-api-key")

	config := stackletProviderModel{
		Endpoint: types.StringNull(),
		APIKey:   types.StringNull(),
	}

	creds := getCredentials(config)
	assert.Equal(t, "https://env-endpoint.example.com", creds.Endpoint)
	assert.Equal(t, "env-api-key", creds.APIKey)
}

func TestGetCredentials_ConfigOverridesEnvVars(t *testing.T) {
	t.Setenv("STACKLET_ENDPOINT", "https://env-endpoint.example.com")
	t.Setenv("STACKLET_API_KEY", "env-api-key")

	config := stackletProviderModel{
		Endpoint: types.StringValue("https://config-endpoint.example.com"),
		APIKey:   types.StringValue("config-api-key"),
	}

	creds := getCredentials(config)
	assert.Equal(t, "https://config-endpoint.example.com", creds.Endpoint)
	assert.Equal(t, "config-api-key", creds.APIKey)
}

func TestGetCredentials_FromStackletAdminConfig(t *testing.T) {
	setupStackletAdminConfig(t, "https://cli-endpoint.example.com", "cli-api-key")

	config := stackletProviderModel{
		Endpoint: types.StringNull(),
		APIKey:   types.StringNull(),
	}

	creds := getCredentials(config)
	assert.Equal(t, "https://cli-endpoint.example.com", creds.Endpoint)
	assert.Equal(t, "cli-api-key", creds.APIKey)
}

func TestGetCredentials_EnvVarsOverrideStackletAdminFiles(t *testing.T) {
	t.Setenv("STACKLET_ENDPOINT", "https://env-endpoint.example.com")
	t.Setenv("STACKLET_API_KEY", "env-api-key")

	setupStackletAdminConfig(t, "https://cli-endpoint.example.com", "cli-api-key")

	config := stackletProviderModel{
		Endpoint: types.StringNull(),
		APIKey:   types.StringNull(),
	}

	creds := getCredentials(config)
	assert.Equal(t, "https://env-endpoint.example.com", creds.Endpoint)
	assert.Equal(t, "env-api-key", creds.APIKey)
}

func TestGetCredentials_MixedSources(t *testing.T) {
	t.Setenv("STACKLET_API_KEY", "env-api-key")

	setupStackletAdminConfig(t, "https://cli-endpoint.example.com", "")

	config := stackletProviderModel{
		Endpoint: types.StringNull(),
		APIKey:   types.StringNull(),
	}

	creds := getCredentials(config)
	assert.Equal(t, "https://cli-endpoint.example.com", creds.Endpoint)
	assert.Equal(t, "env-api-key", creds.APIKey)
}

func TestGetCredentials_Empty(t *testing.T) {
	setupStackletAdminConfig(t, "", "")

	config := stackletProviderModel{
		Endpoint: types.StringNull(),
		APIKey:   types.StringNull(),
	}

	creds := getCredentials(config)

	assert.Equal(t, "", creds.Endpoint)
	assert.Equal(t, "", creds.APIKey)
}

func TestGetCredentials_PartialConfig(t *testing.T) {
	t.Setenv("STACKLET_API_KEY", "env-api-key")

	config := stackletProviderModel{
		Endpoint: types.StringValue("https://config-endpoint.example.com"),
		APIKey:   types.StringNull(),
	}

	creds := getCredentials(config)

	assert.Equal(t, "https://config-endpoint.example.com", creds.Endpoint)
	assert.Equal(t, "env-api-key", creds.APIKey)
}

func TestGetCredentials_InvalidStackletAdminConfigJSON(t *testing.T) {
	stackletDir := makeStackletAdminDir(t)
	configFile := path.Join(stackletDir, "config.json")
	if err := os.WriteFile(configFile, []byte("invalid json"), 0644); err != nil {
		t.Fatal(err)
	}

	config := stackletProviderModel{
		Endpoint: types.StringNull(),
		APIKey:   types.StringNull(),
	}

	creds := getCredentials(config)
	assert.Equal(t, "", creds.Endpoint)
	assert.Equal(t, "", creds.APIKey)
}

func TestGetCredentials_MissingStackletAdminFiles(t *testing.T) {
	setupStackletAdminConfig(t, "", "")

	config := stackletProviderModel{
		Endpoint: types.StringNull(),
		APIKey:   types.StringNull(),
	}

	creds := getCredentials(config)
	assert.Equal(t, "", creds.Endpoint)
	assert.Equal(t, "", creds.APIKey)
}
