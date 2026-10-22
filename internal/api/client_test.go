// Copyright Stacklet, Inc. 2025, 2026

package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflogtest"
	"github.com/stretchr/testify/assert"
)

func TestUserAgentHeader(t *testing.T) {
	version := "1.2.3"
	expectedUserAgent := "terraform-provider-stacklet/" + version

	var capturedUserAgent string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUserAgent = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	transport := &authTransport{
		APIKey:  "test-api-key",
		Version: version,
		Base:    http.DefaultTransport,
	}

	client := &http.Client{Transport: transport}
	req, _ := http.NewRequest("GET", server.URL, nil)
	_, err := client.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, expectedUserAgent, capturedUserAgent)
}

func TestAuthorizationHeader(t *testing.T) {
	apiKey := "test-api-key-123"
	expectedAuth := "Bearer " + apiKey

	var capturedAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	transport := &authTransport{
		APIKey:  apiKey,
		Version: "1.0.0",
		Base:    http.DefaultTransport,
	}

	client := &http.Client{Transport: transport}
	req, _ := http.NewRequest("GET", server.URL, nil)
	_, err := client.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, expectedAuth, capturedAuth)
}

func TestErrorTransport_400WithGraphQLErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"errors":[{"message":"field not found"},{"message":"invalid type"}]}`))
	}))
	defer server.Close()

	transport := &errorTransport{Base: http.DefaultTransport}
	client := &http.Client{Transport: transport}
	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := client.Do(req)

	assert.Nil(t, resp)
	assert.EqualError(t, err, "Get \""+server.URL+"\": field not found\ninvalid type")
}

func TestErrorTransport_400WithoutGraphQLErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`not json`))
	}))
	defer server.Close()

	transport := &errorTransport{Base: http.DefaultTransport}
	client := &http.Client{Transport: transport}
	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := client.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "not json", string(body))
}

func TestMutationMasksSecretsInDebugLogs(t *testing.T) {
	t.Setenv("TF_LOG", "DEBUG")
	apiKey := "plain-secret-api-key"
	// JSON-encoding escapes this one, so the request body carries a form that
	// differs from the value itself.
	appKey := "app\"secret\nkey"

	var reqBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		reqBody = string(body)
		_, _ = w.Write([]byte(`{"data":{"updateDatadogIntegration":{"datadogIntegration":{"configured":true,"statusInfo":{"details":[]}},"problems":[]}}}`))
	}))
	defer server.Close()

	var output bytes.Buffer
	ctx := tflogtest.RootLogger(context.Background(), &output)
	api := New(ClientConfig{Endpoint: server.URL, APIKey: "token"})
	_, err := api.DatadogIntegration.Update(ctx, DatadogIntegrationInput{APIKey: &apiKey, AppKey: &appKey})
	assert.NoError(t, err)
	assert.Contains(t, reqBody, apiKey, "the request itself carries the secret")

	entries, err := tflogtest.MultilineJSONDecode(&output)
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
	for _, entry := range entries {
		logged := fmt.Sprint(entry["req_body"])
		assert.Contains(t, logged, `"apiKey":"***"`)
		assert.Contains(t, logged, `"appKey":"***"`)
		assert.NotContains(t, logged, "secret")
	}
}

func TestErrorTransport_NonBadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{}}`))
	}))
	defer server.Close()

	transport := &errorTransport{Base: http.DefaultTransport}
	client := &http.Client{Transport: transport}
	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := client.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
