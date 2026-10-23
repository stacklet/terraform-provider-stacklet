// Copyright Stacklet, Inc. 2025, 2026

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
	Request  recordedRequest  `json:"request"`
	Response recordedResponse `json:"response"`
}

// recordedRequest is a recorded GraphQL or REST request.
type recordedRequest struct {
	// GraphQL requests
	Query     string         `json:"query,omitempty"`
	Variables map[string]any `json:"variables,omitempty"`

	// REST requests
	Method string          `json:"method,omitempty"`
	Path   string          `json:"path,omitempty"`
	Body   json.RawMessage `json:"body,omitempty"`
}

func (r recordedRequest) IsGraphQL() bool {
	return r.Query != ""
}

func (r recordedRequest) Key() string {
	if r.IsGraphQL() {
		key := r.Query
		if len(r.Variables) > 0 {
			varsJSON, _ := json.Marshal(r.Variables)
			key = fmt.Sprintf("%s:%s", key, string(varsJSON))
		}
		return key
	}

	key := r.Method + " " + r.Path
	if len(r.Body) > 0 {
		// recordings are saved indented, compact for consistent keys
		var body bytes.Buffer
		if err := json.Compact(&body, r.Body); err == nil {
			key = fmt.Sprintf("%s:%s", key, body.String())
		}
	}
	return key
}

// recordedResponse is a recorded GraphQL or REST response.
type recordedResponse struct {
	// GraphQL responses
	Data   any `json:"data,omitempty"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`

	// REST responses
	StatusCode int             `json:"status_code,omitempty"`
	Body       json.RawMessage `json:"body,omitempty"`
}

// graphqlRequest is used to detect GraphQL requests.
type graphqlRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

// graphqlResponse is used to decode GraphQL responses.
type graphqlResponse struct {
	Data   any `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
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
		if err := os.MkdirAll(dirname, 0o755); err != nil {
			return fmt.Errorf("failed to create recordings directory: %v", err)
		}
	}

	absPath, err := filepath.Abs(filename)
	if err != nil {
		return err
	}
	rt.t.Logf("Saving recordings to: %s", absPath)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open recordings file: %v", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(rt.recordings); err != nil {
		return fmt.Errorf("failed to marshal recordings: %v", err)
	}
	return nil
}

func (rt *recordedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// on errors, return an empty response, as returning nil causes
	// http.Client.Do to traceback.
	emptyResp := &http.Response{}

	var body []byte
	if req.Body != nil {
		var err error
		body, err = io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {
			return emptyResp, fmt.Errorf("failed to read request body: %v", err)
		}
	}

	recReq, err := newRecordedRequest(req, body)
	if err != nil {
		return emptyResp, err
	}

	if rt.mode == TestModeRecord {
		return rt.processRecordRequest(req, body, recReq)
	}
	return rt.processReplayRequest(recReq)
}

// newRecordedRequest returns a recordedRequest for a GraphQL or REST request.
func newRecordedRequest(req *http.Request, body []byte) (recordedRequest, error) {
	var gqlReq graphqlRequest
	if len(body) > 0 && json.Unmarshal(body, &gqlReq) == nil && gqlReq.Query != "" {
		return recordedRequest{Query: gqlReq.Query, Variables: gqlReq.Variables}, nil
	}

	recReq := recordedRequest{Method: req.Method, Path: req.URL.RequestURI()}
	if len(body) > 0 {
		if !json.Valid(body) {
			return recordedRequest{}, fmt.Errorf("request body is not valid JSON: %s", body)
		}
		recReq.Body = body
	}
	return recReq, nil
}

func (rt *recordedTransport) processRecordRequest(req *http.Request, body []byte, recReq recordedRequest) (*http.Response, error) {
	// on errors, return an empty response, as returning nil causes
	// http.Client.Do to traceback.
	emptyResp := &http.Response{}

	rt.t.Logf("Recording request: %s", recReq.Key())

	// Make the real request with the original body
	if body != nil {
		req.Body = io.NopCloser(bytes.NewReader(body))
	}
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

	var recResp recordedResponse
	if recReq.IsGraphQL() {
		var gqlResp graphqlResponse
		if err := json.Unmarshal(respBody, &gqlResp); err != nil {
			return emptyResp, fmt.Errorf("failed to decode response body: %v", err)
		}
		recResp.Data = gqlResp.Data
		recResp.Errors = gqlResp.Errors
	} else {
		recResp.StatusCode = resp.StatusCode
		if len(respBody) > 0 {
			if !json.Valid(respBody) {
				return emptyResp, fmt.Errorf("response body is not valid JSON: %s", respBody)
			}
			recResp.Body = respBody
		}
	}

	rt.addRecording(recording{Request: recReq, Response: recResp})

	// Return the response with a new body reader
	resp.Body = io.NopCloser(bytes.NewReader(respBody))
	return resp, nil
}

func (rt *recordedTransport) processReplayRequest(recReq recordedRequest) (*http.Response, error) {
	// on errors, return an empty response, as returning nil causes
	// http.Client.Do to traceback.
	emptyResp := &http.Response{}

	key := recReq.Key()
	rt.t.Logf("Attempting to replay request: %s", key)
	rec, err := rt.getRecording(key)
	if err != nil {
		return emptyResp, err
	}

	if !recReq.IsGraphQL() {
		statusCode := rec.Response.StatusCode
		if statusCode == 0 {
			statusCode = http.StatusOK
		}
		return &http.Response{
			StatusCode: statusCode,
			Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
			Body:       io.NopCloser(bytes.NewReader(rec.Response.Body)),
			Header:     make(http.Header),
		}, nil
	}

	// Check that the recording matches
	if recReq.Query != rec.Request.Query || !reflect.DeepEqual(recReq.Variables, rec.Request.Variables) {
		return emptyResp, fmt.Errorf(`
Request doesn't match expected one:
-- query
got     : %s
expected: %s
-- parameters
got     : %v
expected: %v`,
			recReq.Query, rec.Request.Query, recReq.Variables, rec.Request.Variables)
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
	return recording{}, fmt.Errorf("no recording found for request: %s", key)
}
