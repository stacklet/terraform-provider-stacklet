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
	t          *testing.T
	testName   string
	mode       string // test run mode
	wrapped    http.RoundTripper
	recChans   map[string](chan recording) // channels returning replayed recordings
	recordings map[string][]recording
	recLock    sync.Mutex
}

type recording struct {
	Request  graphqlRequest  `json:"request"`
	Response graphqlResponse `json:"response"`
}

type graphqlRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

func (g graphqlRequest) Key() string {
	key := g.Query
	if len(g.Variables) > 0 {
		varsJSON, _ := json.Marshal(g.Variables)
		key = fmt.Sprintf("%s:%s", key, string(varsJSON))
	}
	return key
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
		t:          t,
		testName:   testName,
		mode:       mode,
		wrapped:    wrapped,
		recChans:   make(map[string](chan recording)),
		recordings: make(map[string][]recording),
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

	recordings := make(map[string][]recording)
	if err := json.Unmarshal(data, &recordings); err != nil {
		return err
	}
	for key, recs := range recordings {
		ch := make(chan recording, len(recs))
		for _, rec := range recs {
			ch <- rec
		}
		rt.recChans[key] = ch
	}
	return nil
}

func (rt *recordedTransport) saveRecording() error {
	if rt.mode != TestModeRecord {
		return nil
	}

	dirname := filepath.Join("recordings")
	filename := filepath.Join(dirname, rt.testName+".json")

	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		rt.t.Logf("Creating recordings directory: %s", dirname)
		if err := os.MkdirAll(dirname, 0755); err != nil {
			return fmt.Errorf("failed to create recordings directory: %v", err)
		}
	}

	absPath, err := filepath.Abs(filename)
	if err != nil {
		return err
	}
	rt.t.Logf("Saving recordings to: %s", absPath)
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
	// on errors, return an empty response, as returning nil causes
	// http.Client.Do to traceback.
	emptyResp := &http.Response{}

	// Read the request body
	var gqlReq graphqlRequest
	body, err := io.ReadAll(req.Body)
	if err != nil {
		req.Body.Close()
		return emptyResp, fmt.Errorf("failed to read request body: %v", err)
	}
	if err := json.Unmarshal(body, &gqlReq); err != nil {
		req.Body.Close()
		return emptyResp, fmt.Errorf("failed to decode request body: %v", err)
	}

	if rt.mode == TestModeRecord {
		return rt.processRecordRequest(req, body, gqlReq)
	}
	return rt.processReplayRequest(gqlReq)
}

func (rt *recordedTransport) processRecordRequest(req *http.Request, body []byte, gqlReq graphqlRequest) (*http.Response, error) {
	// on errors, return an empty response, as returning nil causes
	// http.Client.Do to traceback.
	emptyResp := &http.Response{}

	rt.t.Logf("Recording request with query: %s - %v", gqlReq.Query, gqlReq.Variables)

	// Make the real request with the original body
	req.Body = io.NopCloser(bytes.NewReader(body))
	resp, err := rt.wrapped.RoundTrip(req)
	if err != nil {
		return emptyResp, err
	}

	// Read and store the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return emptyResp, fmt.Errorf("failed to read response body: %v", err)
	}
	if err := resp.Body.Close(); err != nil {
		return emptyResp, err
	}

	// Parse the response
	var gqlResp graphqlResponse
	if err := json.Unmarshal(respBody, &gqlResp); err != nil {
		return emptyResp, fmt.Errorf("failed to decode response body: %v", err)
	}

	rt.addRecording(recording{Request: gqlReq, Response: gqlResp})

	// Return the response with a new body reader
	resp.Body = io.NopCloser(bytes.NewReader(respBody))
	return resp, nil
}

func (rt *recordedTransport) processReplayRequest(gqlReq graphqlRequest) (*http.Response, error) {
	// on errors, return an empty response, as returning nil causes
	// http.Client.Do to traceback.
	emptyResp := &http.Response{}

	key := gqlReq.Key()
	rt.t.Logf("Attempting to replay request with query: %s - %v", gqlReq.Query, gqlReq.Variables)
	rec, err := rt.getRecording(key)
	if err != nil {
		return emptyResp, err
	}

	// Check that the recording matches
	if gqlReq.Query != rec.Request.Query || !reflect.DeepEqual(gqlReq.Variables, rec.Request.Variables) {
		return emptyResp, fmt.Errorf(`
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
		return emptyResp, fmt.Errorf("failed to marshal recorded response: %v", err)
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(respBody)),
		Header:     make(http.Header),
	}, nil
}

func (rt *recordedTransport) addRecording(r recording) {
	rt.recLock.Lock()
	defer rt.recLock.Unlock()
	key := r.Request.Key()
	rt.recordings[key] = append(rt.recordings[key], r)
}

func (rt *recordedTransport) getRecording(key string) (recording, error) {
	if ch, ok := rt.recChans[key]; ok {
		select {
		case rec := <-ch:
			rt.t.Logf("Recording found with key: %s", key)
			return rec, nil
		default:
			// no recording found in channel, fall out to the error case
		}
	}
	return recording{}, fmt.Errorf("no recording found for query: %s", key)
}
