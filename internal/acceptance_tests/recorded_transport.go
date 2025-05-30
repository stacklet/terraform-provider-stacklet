// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
)

type recordedTransport struct {
	recordings     map[string][]recording
	recordingsLock sync.Mutex // lock around modifications of recordings
	mode           string     // test run mode
	t              *testing.T
	testName       string
	wrapped        http.RoundTripper
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

func newRecordedTransport(t *testing.T, testName string, mode string, wrapped http.RoundTripper) *recordedTransport {
	t.Logf("%s mode enabled for test: %s", mode, testName)

	return &recordedTransport{
		recordings: make(map[string][]recording),
		mode:       mode,
		t:          t,
		testName:   testName,
		wrapped:    wrapped,
	}
}

func (rt *recordedTransport) loadRecording() error {
	if rt.mode != TestModeReplay {
		return nil
	}

	filename := filepath.Join("recordings", rt.testName+".json")
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read recordings: %v", err)
	}

	return json.Unmarshal(data, &rt.recordings)
}

func (rt *recordedTransport) saveRecording() error {
	if rt.mode != TestModeRecord {
		return nil
	}

	dirname := filepath.Join("recordings")
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
		req.Body.Close()
		return nil, fmt.Errorf("failed to read request body: %v", err)
	}
	req.Body = io.NopCloser(bytes.NewReader(body)) // Reset the body for future reads

	if err := json.Unmarshal(body, &gqlReq); err != nil {
		req.Body.Close()
		return nil, fmt.Errorf("failed to decode request body: %v", err)
	}

	// Create a key that includes both the query and variables
	key := gqlReq.Query
	if len(gqlReq.Variables) > 0 {
		varsJSON, _ := json.Marshal(gqlReq.Variables)
		key = fmt.Sprintf("%s:%s", key, string(varsJSON))
	}

	if rt.mode == TestModeRecord {
		rt.t.Logf("Recording request with query: %s - %v", gqlReq.Query, gqlReq.Variables)

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
		rt.recordingsLock.Lock()
		rt.recordings[key] = append(rt.recordings[key], recording{
			Request:  gqlReq,
			Response: gqlResp,
		})
		rt.recordingsLock.Unlock()

		// Return the response with a new body reader
		resp.Body = io.NopCloser(bytes.NewReader(respBody))
		return resp, nil
	}

	// Replay mode
	rt.t.Logf("Attempting to replay request with query: %s - %v", gqlReq.Query, gqlReq.Variables)
	rt.recordingsLock.Lock()
	recs, ok := rt.recordings[key]
	if !ok || len(recs) == 0 {
		return nil, fmt.Errorf("no recording found for query: %s", key)
	}
	rt.t.Logf("Recording found with key: %s", key)

	// Use the first recording and rotate
	rec := recs[0]
	rt.recordings[key] = recs[1:]
	rt.recordingsLock.Unlock()

	// Check that the recording matches
	if gqlReq.Query != rec.Request.Query || !reflect.DeepEqual(gqlReq.Variables, rec.Request.Variables) {
		return nil, fmt.Errorf(`
Request doesn't match expected one:
-- query
got     : %s
expected: %s
-- parameters
got     : %v
expected: %v`,
			gqlReq.Query, rec.Request.Query, gqlReq.Variables, rec.Request.Variables)
	}

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
