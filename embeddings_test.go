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

func TestEmbeddings(t *testing.T) {
	req, _ := New("hunter2")
	req.Client.RetryMax = 1
	req.BaseURL = ""
	req.Debug = false

	// Case: Embeddings Fails with no context
	resp, err := req.Embeddings(nil, "", "", EmbeddingsRequest{}) //lint:ignore SA1012 nil context used intentionally
	if !reflect.DeepEqual(resp, EmbeddingsResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, EmbeddingsResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "no context provided")
	}

	// Case: Completions Fails with no model
	resp, err = req.Embeddings(context.TODO(), "", "", EmbeddingsRequest{})
	if !reflect.DeepEqual(resp, EmbeddingsResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, EmbeddingsResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "no model provided")
	}

	// Case: Completions Fails with no input
	resp, err = req.Embeddings(context.TODO(), "a", "", EmbeddingsRequest{})
	if !reflect.DeepEqual(resp, EmbeddingsResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, EmbeddingsResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "no prompt provided")
	}

	// Case: Chat Completion Fails with invalid HTTP request
	resp, err = req.Embeddings(context.TODO(), "a", "b", EmbeddingsRequest{})
	if !reflect.DeepEqual(resp, EmbeddingsResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, EmbeddingsResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "!nil")
	}

	// Case: Chat Completion Fails with HTTP 200 and malformed body
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))

	req.BaseURL = ts.URL

	resp, err = req.Embeddings(context.TODO(), "a", "b", EmbeddingsRequest{})
	if !reflect.DeepEqual(resp, EmbeddingsResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, EmbeddingsResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "!nil")
	}

	ts.Close()

	// Case: Chat Completion Fails with HTTP 400
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprintln(w, nil)
	}))

	req.BaseURL = ts.URL

	resp, err = req.Embeddings(context.TODO(), "a", "b", EmbeddingsRequest{})
	if !reflect.DeepEqual(resp, EmbeddingsResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, EmbeddingsResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "!nil")
	}

	ts.Close()

	// Case: Chat Completion Succeeds with HTTP 200
	req.Debug = true
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, _ := json.Marshal(EmbeddingsResponse{})
		fmt.Fprintln(w, bytes.NewBuffer(resp))
	}))

	req.BaseURL = ts.URL

	resp, err = req.Embeddings(context.TODO(), "a", "b", EmbeddingsRequest{})
	if !reflect.DeepEqual(resp, EmbeddingsResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, EmbeddingsResponse{})
	}
	if err != nil {
		t.Errorf("Error was incorrect, got: %s, want: %v.", err, nil)
	}

	ts.Close()
}
