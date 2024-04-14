package together

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

// Test creating a new API object/client with a valid API key.
func TestNew(t *testing.T) {
	result, _ := New("hunter2")
	if result.APIKey != "hunter2" {
		t.Errorf("Result was incorrect, got: %s, want: %s.", result.APIKey, "hunter2")
	}

	result, err := New("")
	if result != nil {
		t.Errorf("Result was incorrect, got: %s, want: %s.", result.APIKey, "")
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "invalid credentials: API Token must not be empty")
	}

}

// Test copying headers from one header map to another.
func TestCopyHeader(t *testing.T) {
	a := make(http.Header)
	a.Set("test", "test")

	b := make(http.Header)
	copyHeader(b, a)

	if len(b) == 0 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", len(b), 1)
	}

	if b.Get("test") != "test" {
		t.Errorf("Result was incorrect, got: %s, want: %s.", b.Get("test"), "test")
	}
}

// Test making a HTTP request.
func TestRequest(t *testing.T) {
	req, _ := New("hunter2")
	req.Client.RetryMax = 1
	req.BaseURL = ""

	// Case: HTTP request creation fails
	resp, err := req.request(nil, "", "", nil, nil) //lint:ignore SA1012 nil context used intentionally
	if err == nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if resp != nil {
		t.Errorf("Result was incorrect, got: %v, want: non-nil", resp)
	}

	// Case: HTTP request fails
	req.BaseURL = "https://test.test" // RFC2606§2 reserved TLD
	resp, err = req.request(context.Background(), "GET", "", nil, nil)
	if err == nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if resp != nil {
		t.Errorf("Result was incorrect, got: %v, want: non-nil", resp)
	}

	// Case: HTTP request succeeds
	req.Debug = true
	req.BaseURL = "https://postman-echo.com/post"

	reqBody, _ := json.Marshal(map[string]string{"key": "value"})

	resp, err = req.request(context.Background(), "POST", "", bytes.NewBuffer(reqBody), nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if resp == nil {
		t.Errorf("Result was incorrect, got: %v, want: non-nil", resp)
	}

}
