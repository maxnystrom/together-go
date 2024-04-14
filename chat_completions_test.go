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

func TestChatCompletions(t *testing.T) {
	req, _ := New("hunter2")
	req.Client.RetryMax = 1
	req.BaseURL = ""
	req.Debug = false

	// Case: Chat Completion Fails with no message
	resp, err := req.ChatCompletions(context.TODO(), "", nil, ChatCompletionsRequest{})
	if !reflect.DeepEqual(resp, ChatCompletionsResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, ChatCompletionsResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "no messages provided")
	}

	// Case: Chat Completion Fails with no model
	resp, err = req.ChatCompletions(context.TODO(), "", []Message{{"Role", "Content"}}, ChatCompletionsRequest{})
	if !reflect.DeepEqual(resp, ChatCompletionsResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, ChatCompletionsResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "no model provided")
	}

	// Case: Chat Completion Fails with no context provided
	resp, err = req.ChatCompletions(nil, "a", []Message{{"Role", "Content"}}, ChatCompletionsRequest{}) //lint:ignore SA1012 nil context used intentionally
	if !reflect.DeepEqual(resp, ChatCompletionsResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, ChatCompletionsResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "no context provided")
	}

	// Case: Chat Completion Fails with invalid HTTP request
	resp, err = req.ChatCompletions(context.TODO(), "a", []Message{{"Role", "Content"}}, ChatCompletionsRequest{})
	if !reflect.DeepEqual(resp, ChatCompletionsResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, ChatCompletionsResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %v.", err, nil)
	}

	// Case: Chat Completion Fails with HTTP 200 and malformed body
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))

	req.BaseURL = ts.URL

	resp, err = req.ChatCompletions(context.TODO(), "a", []Message{{"Role", "Content"}}, ChatCompletionsRequest{})
	if !reflect.DeepEqual(resp, ChatCompletionsResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, ChatCompletionsResponse{})
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

	resp, err = req.ChatCompletions(context.TODO(), "a", []Message{{"Role", "Content"}}, ChatCompletionsRequest{})
	if !reflect.DeepEqual(resp, ChatCompletionsResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, ChatCompletionsResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %v.", err, nil)
	}

	ts.Close()

	// Case: Chat Completion Succeeds with HTTP 200
	req.Debug = true
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, _ := json.Marshal(ChatCompletionsResponse{})
		fmt.Fprintln(w, bytes.NewBuffer(resp))
	}))

	req.BaseURL = ts.URL

	resp, err = req.ChatCompletions(context.TODO(), "a", []Message{{"Role", "Content"}}, ChatCompletionsRequest{})
	if !reflect.DeepEqual(resp, ChatCompletionsResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, ChatCompletionsResponse{})
	}
	if err != nil {
		t.Errorf("Error was incorrect, got: %s, want: %v.", err, nil)
	}

	ts.Close()
}
