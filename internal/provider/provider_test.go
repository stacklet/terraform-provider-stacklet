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
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	stackletDir := path.Join(tmpDir, ".stacklet")
	if err := os.MkdirAll(stackletDir, 0755); err != nil {
		t.Fatal(err)
	}

	configFile := path.Join(stackletDir, "config.json")
	configData := map[string]string{"api": "https://cli-endpoint.example.com"}
	configJSON, _ := json.Marshal(configData)
	if err := os.WriteFile(configFile, configJSON, 0644); err != nil {
		t.Fatal(err)
	}

	credsFile := path.Join(stackletDir, "credentials")
	if err := os.WriteFile(credsFile, []byte("cli-api-key"), 0644); err != nil {
		t.Fatal(err)
	}

	config := stackletProviderModel{
		Endpoint: types.StringNull(),
		APIKey:   types.StringNull(),
	}

	creds := getCredentials(config)

	assert.Equal(t, "https://cli-endpoint.example.com", creds.Endpoint)
	assert.Equal(t, "cli-api-key", creds.APIKey)
}

func TestGetCredentials_EnvVarsOverrideStackletAdminFiles(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("STACKLET_ENDPOINT", "https://env-endpoint.example.com")
	t.Setenv("STACKLET_API_KEY", "env-api-key")

	stackletDir := path.Join(tmpDir, ".stacklet")
	if err := os.MkdirAll(stackletDir, 0755); err != nil {
		t.Fatal(err)
	}

	configFile := path.Join(stackletDir, "config.json")
	configData := map[string]string{"api": "https://cli-endpoint.example.com"}
	configJSON, _ := json.Marshal(configData)
	if err := os.WriteFile(configFile, configJSON, 0644); err != nil {
		t.Fatal(err)
	}

	credsFile := path.Join(stackletDir, "credentials")
	if err := os.WriteFile(credsFile, []byte("cli-api-key"), 0644); err != nil {
		t.Fatal(err)
	}

	config := stackletProviderModel{
		Endpoint: types.StringNull(),
		APIKey:   types.StringNull(),
	}

	creds := getCredentials(config)

	assert.Equal(t, "https://env-endpoint.example.com", creds.Endpoint)
	assert.Equal(t, "env-api-key", creds.APIKey)
}

func TestGetCredentials_MixedSources(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("STACKLET_API_KEY", "env-api-key")

	stackletDir := path.Join(tmpDir, ".stacklet")
	if err := os.MkdirAll(stackletDir, 0755); err != nil {
		t.Fatal(err)
	}

	configFile := path.Join(stackletDir, "config.json")
	configData := map[string]string{"api": "https://cli-endpoint.example.com"}
	configJSON, _ := json.Marshal(configData)
	if err := os.WriteFile(configFile, configJSON, 0644); err != nil {
		t.Fatal(err)
	}

	config := stackletProviderModel{
		Endpoint: types.StringNull(),
		APIKey:   types.StringNull(),
	}

	creds := getCredentials(config)

	assert.Equal(t, "https://cli-endpoint.example.com", creds.Endpoint)
	assert.Equal(t, "env-api-key", creds.APIKey)
}

func TestGetCredentials_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

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
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	stackletDir := path.Join(tmpDir, ".stacklet")
	if err := os.MkdirAll(stackletDir, 0755); err != nil {
		t.Fatal(err)
	}

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
}

func TestGetCredentials_MissingStackletAdminFiles(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	stackletDir := path.Join(tmpDir, ".stacklet")
	if err := os.MkdirAll(stackletDir, 0755); err != nil {
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
