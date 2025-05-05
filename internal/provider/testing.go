package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

type recordedTransport struct {
	recordings map[string][]recording
	mode       string // "record" or "replay"
	t          *testing.T
	testName   string
	wrapped    http.RoundTripper
}

type recording struct {
	Request  graphqlRequest  `json:"request"`
	Response graphqlResponse `json:"response"`
}

type graphqlRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type graphqlResponse struct {
	Data   any `json:"data,omitempty"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

func newRecordedTransport(t *testing.T, testName string, wrapped http.RoundTripper) *recordedTransport {
	mode := "replay"
	if os.Getenv("STACKLET_RECORD") != "" {
		mode = "record"
		t.Logf("Recording mode enabled for test: %s", testName)
	} else {
		t.Logf("Replay mode enabled for test: %s", testName)
	}

	return &recordedTransport{
		recordings: make(map[string][]recording),
		mode:       mode,
		t:          t,
		testName:   testName,
		wrapped:    wrapped,
	}
}

func (rt *recordedTransport) loadRecordings() error {
	if rt.mode == "record" {
		return nil
	}

	filename := filepath.Join("testdata", "recordings", rt.testName+".json")
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read recordings: %v", err)
	}

	return json.Unmarshal(data, &rt.recordings)
}

func (rt *recordedTransport) saveRecordings() error {
	if rt.mode != "record" {
		return nil
	}

	dirname := filepath.Join("testdata", "recordings")
	filename := filepath.Join(dirname, rt.testName+".json")

	// Check if directory exists
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		rt.t.Logf("Creating recordings directory: %s", dirname)
		if err := os.MkdirAll(dirname, 0755); err != nil {
			return fmt.Errorf("failed to create recordings directory: %v", err)
		}
	}

	// Log the current working directory and absolute path
	cwd, err := os.Getwd()
	if err == nil {
		rt.t.Logf("Current working directory: %s", cwd)
		absPath, err := filepath.Abs(filename)
		if err == nil {
			rt.t.Logf("Saving recordings to: %s", absPath)
		}
	}

	data, err := json.MarshalIndent(rt.recordings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal recordings: %v", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write recordings file: %v", err)
	}

	return nil
}

func (rt *recordedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Read the request body
	var gqlReq graphqlRequest
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %v", err)
	}
	req.Body = io.NopCloser(bytes.NewReader(body)) // Reset the body for future reads

	if err := json.Unmarshal(body, &gqlReq); err != nil {
		return nil, fmt.Errorf("failed to decode request body: %v", err)
	}

	// Create a key that includes both the query and variables
	key := gqlReq.Query
	if len(gqlReq.Variables) > 0 {
		varsJSON, _ := json.Marshal(gqlReq.Variables)
		key = fmt.Sprintf("%s:%s", key, string(varsJSON))
	}

	if rt.mode == "record" {
		rt.t.Logf("Recording request with query: %s", gqlReq.Query)

		// Make the real request with the original body
		req.Body = io.NopCloser(bytes.NewReader(body))
		resp, err := rt.wrapped.RoundTrip(req)
		if err != nil {
			return nil, err
		}

		// Read and store the response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}
		if err := resp.Body.Close(); err != nil {
			return nil, err
		}

		// Parse the response
		var gqlResp graphqlResponse
		if err := json.Unmarshal(respBody, &gqlResp); err != nil {
			return nil, fmt.Errorf("failed to decode response body: %v", err)
		}

		// Record the interaction
		rt.recordings[key] = append(rt.recordings[key], recording{
			Request:  gqlReq,
			Response: gqlResp,
		})

		// Save recordings after each request
		if err := rt.saveRecordings(); err != nil {
			rt.t.Logf("Warning: failed to save recordings: %v", err)
		} else {
			rt.t.Logf("Successfully saved recording to testdata/recordings/%s.json", rt.testName)
		}

		// Return the response with a new body reader
		resp.Body = io.NopCloser(bytes.NewReader(respBody))
		return resp, nil
	}

	// Replay mode
	rt.t.Logf("Attempting to replay request with query: %s", gqlReq.Query)
	recs, ok := rt.recordings[key]
	if !ok || len(recs) == 0 {
		return nil, fmt.Errorf("no recording found for query: %s", key)
	}

	// Use the first recording and rotate
	rec := recs[0]
	rt.recordings[key] = recs[1:]

	// Create response
	respBody, err := json.Marshal(rec.Response)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal recorded response: %v", err)
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(respBody)),
		Header:     make(http.Header),
	}, nil
}

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a new provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"stacklet": func() (tfprotov6.ProviderServer, error) {
		p := New("test")()
		return providerserver.NewProtocol6WithError(p)()
	},
}

func testAccPreCheck(t *testing.T) {
	// Verify environment variables required for API calls are set
	if v := os.Getenv("STACKLET_API_KEY"); v == "" {
		t.Skip("STACKLET_API_KEY must be set for acceptance tests")
	}
}

// setupRecordedTest prepares a test with recording/replay capability
func setupRecordedTest(t *testing.T, testName string) (*recordedTransport, error) {
	rt := newRecordedTransport(t, testName, http.DefaultTransport)
	if err := rt.loadRecordings(); err != nil && rt.mode == "replay" {
		return nil, fmt.Errorf("failed to load recordings: %v", err)
	}

	t.Cleanup(func() {
		if err := rt.saveRecordings(); err != nil {
			t.Errorf("failed to save recordings: %v", err)
		}
	})

	return rt, nil
}

func testAccCheckMapValues(resourceName string, mapAttr string, expectedMap map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		for k, v := range expectedMap {
			attrKey := fmt.Sprintf("%s.%s", mapAttr, k)
			if rs.Primary.Attributes[attrKey] != v {
				return fmt.Errorf("expected %s to be %s, got %s", attrKey, v, rs.Primary.Attributes[attrKey])
			}
		}

		return nil
	}
}
