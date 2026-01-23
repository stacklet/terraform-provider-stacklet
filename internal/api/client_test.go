// Copyright (c) 2025 - Stacklet, Inc.

package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
