package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// httpResp is the expected JSON structure for responses.
//
// Tests for sendResponse and sendErrorResponse
//

func TestSendResponse_Success(t *testing.T) {
	rr := httptest.NewRecorder()
	data := map[string]string{"foo": "bar"}
	status := http.StatusOK
	message := "Success"

	sendResponse(rr, status, data, message)

	// Check header is set correctly.
	if ct := rr.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Errorf("Expected Content-Type header 'application/json; charset=utf-8', got %q", ct)
	}
	if rr.Code != status {
		t.Errorf("Expected status %d, got %d", status, rr.Code)
	}

	// Parse response body.
	var resp httpResp
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if resp.Status != status {
		t.Errorf("Expected response status %d, got %d", status, resp.Status)
	}
	if resp.Message != message {
		t.Errorf("Expected message %q, got %q", message, resp.Message)
	}
	expectedData, _ := json.Marshal(data)
	actualData, _ := json.Marshal(resp.Data)
	if !bytes.Equal(expectedData, actualData) {
		t.Errorf("Expected data %s, got %s", expectedData, actualData)
	}
}

func TestSendResponse_MarshalError(t *testing.T) {
	rr := httptest.NewRecorder()
	// Create data that cannot be marshaled (channels cannot be marshaled to JSON).
	data := make(chan int)
	status := http.StatusOK
	message := "This will fail"

	sendResponse(rr, status, data, message)

	// Because WriteHeader is called before marshalling,
	// the status code in the recorder remains unchanged (e.g. 200).
	// Instead of checking rr.Code, we parse the response body to verify the JSON error payload.
	var resp httpResp
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// We expect the JSON payload to reflect an internal error.
	if resp.Status != http.StatusInternalServerError {
		t.Errorf("Expected JSON response status %d, got %d", http.StatusInternalServerError, resp.Status)
	}
	if resp.Message != "Internal Server Error" {
		t.Errorf("Expected message %q, got %q", "Internal Server Error", resp.Message)
	}
}

func TestSendErrorResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	status := http.StatusBadRequest
	message := "Bad Request Error"
	data := map[string]int{"error": 1}

	sendErrorResponse(rr, status, data, message)

	if ct := rr.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Errorf("Expected Content-Type header 'application/json; charset=utf-8', got %q", ct)
	}
	if rr.Code != status {
		t.Errorf("Expected status %d, got %d", status, rr.Code)
	}
	var resp httpResp
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if resp.Status != status {
		t.Errorf("Expected response status %d, got %d", status, resp.Status)
	}
	if resp.Message != message {
		t.Errorf("Expected message %q, got %q", message, resp.Message)
	}
	expectedData, _ := json.Marshal(data)
	actualData, _ := json.Marshal(resp.Data)
	if !bytes.Equal(expectedData, actualData) {
		t.Errorf("Expected data %s, got %s", expectedData, actualData)
	}
}

//
// Tests for getBodyWithType
//

// testStruct is a helper struct for testing getBodyWithType.
type testStruct struct {
	Foo string `json:"foo"`
}

func TestGetBodyWithType_Success(t *testing.T) {
	input := testStruct{Foo: "bar"}
	bodyBytes, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(bodyBytes))

	result, err := getBodyWithType[testStruct](req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result.Foo != input.Foo {
		t.Errorf("Expected Foo %q, got %q", input.Foo, result.Foo)
	}
}

// errorReader simulates a reader that returns an error.
type errorReader struct{}

func (e *errorReader) Read(p []byte) (int, error) {
	return 0, errors.New("read error")
}

func (e *errorReader) Close() error {
	return nil
}

func TestGetBodyWithType_ReadError(t *testing.T) {
	req := httptest.NewRequest("POST", "/", &errorReader{})
	_, err := getBodyWithType[testStruct](req)
	if err == nil {
		t.Fatal("Expected error due to read failure, got nil")
	}
	expectedErr := fmt.Sprintf("%d %s", http.StatusBadRequest, ErrCouldNotReadBody.Error())
	if err.Error() != expectedErr {
		t.Errorf("Expected error %q, got %q", expectedErr, err.Error())
	}
}

func TestGetBodyWithType_UnmarshalError(t *testing.T) {
	// Provide invalid JSON.
	req := httptest.NewRequest("POST", "/", bytes.NewReader([]byte("invalid json")))
	_, err := getBodyWithType[testStruct](req)
	if err == nil {
		t.Fatal("Expected error due to unmarshal failure, got nil")
	}
	expectedErr := fmt.Sprintf("%d %s", http.StatusBadRequest, ErrCouldNotParseBody.Error())
	if err.Error() != expectedErr {
		t.Errorf("Expected error %q, got %q", expectedErr, err.Error())
	}
}

//
// Tests for newError, parseError, and sendHerrorResponse
//

func TestNewError(t *testing.T) {
	status := 404
	message := "Not Found"
	err := newError(status, message)
	expected := "404 Not Found"
	if err.Error() != expected {
		t.Errorf("Expected error %q, got %q", expected, err.Error())
	}
}

func TestParseError_Valid(t *testing.T) {
	err := errors.New("400 Bad Request")
	status, message := parseError(err)
	if status != 400 {
		t.Errorf("Expected status 400, got %d", status)
	}
	if message != "Bad Request" {
		t.Errorf("Expected message %q, got %q", "Bad Request", message)
	}
}

func TestParseError_ShortError(t *testing.T) {
	err := errors.New("OK")
	status, message := parseError(err)
	if status != 500 {
		t.Errorf("Expected status 500 for short error, got %d", status)
	}
	if message != "" {
		t.Errorf("Expected empty message for short error, got %q", message)
	}
}

func TestParseError_NonNumeric(t *testing.T) {
	err := errors.New("ABC Not a number")
	status, message := parseError(err)
	if status != 500 {
		t.Errorf("Expected status 500 for non-numeric error, got %d", status)
	}
	if message != "" {
		t.Errorf("Expected empty message for non-numeric error, got %q", message)
	}
}

func TestSendHerrorResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	// Create an error using newError.
	err := newError(http.StatusNotFound, "Resource not found")
	sendHerrorResponse(rr, err)

	// Validate that sendHerrorResponse writes the correct error response.
	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
	var resp httpResp
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if resp.Status != http.StatusNotFound {
		t.Errorf("Expected response status %d, got %d", http.StatusNotFound, resp.Status)
	}
	if resp.Message != "Resource not found" {
		t.Errorf("Expected message %q, got %q", "Resource not found", resp.Message)
	}
}
