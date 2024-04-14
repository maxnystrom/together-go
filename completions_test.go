package together

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCompletions(t *testing.T) {
	req, _ := New("hunter2")
	req.Client.RetryMax = 1
	req.BaseURL = ""
	req.Debug = false

	// Case: Completions Fails with no context
	resp, err := req.Completions(nil, "", "", 0, CompletionsRequest{}) //lint:ignore SA1012 nil context used intentionally
	if !reflect.DeepEqual(resp, CompletionsResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, CompletionsResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "no context provided")
	}

	// Case: Completions Fails with no model
	resp, err = req.Completions(context.TODO(), "", "", 0, CompletionsRequest{})
	if !reflect.DeepEqual(resp, CompletionsResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, CompletionsResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "no model provided")
	}

	// Case: Completions Fails with no prompt
	resp, err = req.Completions(context.TODO(), "a", "", 0, CompletionsRequest{})
	if !reflect.DeepEqual(resp, CompletionsResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, CompletionsResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "no prompt provided")
	}

	// Case: Completions Fails with no tokens
	resp, err = req.Completions(context.TODO(), "a", "b", 0, CompletionsRequest{})
	if !reflect.DeepEqual(resp, CompletionsResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, CompletionsResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "maxTokens must be greater than 0 and less than 2147483647")
	}

	// Case: Chat Completion Fails with invalid HTTP request
	resp, err = req.Completions(context.TODO(), "a", "b", 10, CompletionsRequest{})
	if !reflect.DeepEqual(resp, CompletionsResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, CompletionsResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %v.", err, nil)
	}

	// Case: Chat Completion Fails with HTTP 200 and malformed body
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))

	req.BaseURL = ts.URL

	resp, err = req.Completions(context.TODO(), "a", "b", 10, CompletionsRequest{})
	if !reflect.DeepEqual(resp, CompletionsResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, CompletionsResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %v.", err, nil)
	}

	ts.Close()

	// Case: Chat Completion Fails with HTTP 400
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprintln(w, nil)
	}))

	req.BaseURL = ts.URL

	resp, err = req.Completions(context.TODO(), "a", "b", 10, CompletionsRequest{})
	if !reflect.DeepEqual(resp, CompletionsResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, ChatCompletionsResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %v.", err, nil)
	}

	ts.Close()

	// Case: Chat Completion Succeeds with HTTP 200
	req.Debug = true
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, _ := json.Marshal(CompletionsResponse{})
		fmt.Fprintln(w, bytes.NewBuffer(resp))
	}))

	req.BaseURL = ts.URL

	resp, err = req.Completions(context.TODO(), "a", "b", 10, CompletionsRequest{})
	if !reflect.DeepEqual(resp, CompletionsResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, CompletionsResponse{})
	}
	if err != nil {
		t.Errorf("Error was incorrect, got: %s, want: %v.", err, nil)
	}

	ts.Close()

}
