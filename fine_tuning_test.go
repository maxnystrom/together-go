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

func TestListRunningInstances(t *testing.T) {
	req, _ := New("hunter2")
	req.Client.RetryMax = 1
	req.BaseURL = ""
	req.Debug = false

	// ListRunningInstances
	//
	// Case: ListRunningInstances Fails with no context
	resp, err := req.ListRunningInstances(nil) //lint:ignore SA1012 nil context used intentionally
	if !reflect.DeepEqual(resp, FineTuningResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, FineTuningResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "no context provided")
	}

	// Case: ListRunningInstances Fails with invalid HTTP request
	resp, err = req.ListRunningInstances(context.TODO())
	if !reflect.DeepEqual(resp, FineTuningResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, FineTuningResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "!nil")
	}

	// Case: ListRunningInstances Fails with HTTP 200 and malformed body
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))

	req.BaseURL = ts.URL

	resp, err = req.ListRunningInstances(context.TODO())
	if !reflect.DeepEqual(resp, FineTuningResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, FineTuningResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "!nil")
	}

	ts.Close()

	// Case: ListRunningInstances Fails with HTTP 400
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprintln(w, nil)
	}))

	req.BaseURL = ts.URL

	resp, err = req.ListRunningInstances(context.TODO())
	if !reflect.DeepEqual(resp, FineTuningResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, FineTuningResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "!nil")
	}

	ts.Close()

	// Case: ListRunningInstances Succeeds with HTTP 200
	req.Debug = true
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, _ := json.Marshal(EmbeddingsResponse{})
		fmt.Fprintln(w, bytes.NewBuffer(resp))
	}))

	req.BaseURL = ts.URL

	resp, err = req.ListRunningInstances(context.TODO())
	if !reflect.DeepEqual(resp, FineTuningResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, FineTuningResponse{})
	}
	if err != nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "!nil")
	}

	ts.Close()

	// StartFineTunedInstance
	//
	// Case: StartFineTunedInstance Fails with no context
	resp, err = req.StartFineTunedInstance(nil, "") //lint:ignore SA1012 nil context used intentionally
	if !reflect.DeepEqual(resp, FineTuningResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, FineTuningResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "no context provided")
	}

	// Case: StartFineTunedInstance Fails with no name
	resp, err = req.StartFineTunedInstance(context.TODO(), "")
	if !reflect.DeepEqual(resp, FineTuningResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, FineTuningResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "no name provided")
	}

	// Case: StartFineTunedInstance Fails with invalid HTTP request
	resp, err = req.StartFineTunedInstance(context.TODO(), "a")
	if !reflect.DeepEqual(resp, FineTuningResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, FineTuningResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "!nil")
	}

	// Case: StartFineTunedInstance Fails with HTTP 200 and malformed body
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))

	req.BaseURL = ts.URL

	resp, err = req.StartFineTunedInstance(context.TODO(), "a")
	if !reflect.DeepEqual(resp, FineTuningResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, FineTuningResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "!nil")
	}

	ts.Close()

	// Case: StartFineTunedInstance Fails with HTTP 400
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprintln(w, nil)
	}))

	req.BaseURL = ts.URL

	resp, err = req.StartFineTunedInstance(context.TODO(), "a")
	if !reflect.DeepEqual(resp, FineTuningResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, FineTuningResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "!nil")
	}

	ts.Close()

	// Case: StartFineTunedInstance Succeeds with HTTP 200
	req.Debug = true
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, _ := json.Marshal(EmbeddingsResponse{})
		fmt.Fprintln(w, bytes.NewBuffer(resp))
	}))

	req.BaseURL = ts.URL

	resp, err = req.StartFineTunedInstance(context.TODO(), "a")
	if !reflect.DeepEqual(resp, FineTuningResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, FineTuningResponse{})
	}
	if err != nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "!nil")
	}

	ts.Close()

	// StopFineTunedInstance
	//
	// Case: StopFineTunedInstance Fails with no context
	resp, err = req.StopFineTunedInstance(nil, "") //lint:ignore SA1012 nil context used intentionally
	if !reflect.DeepEqual(resp, FineTuningResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, FineTuningResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "no context provided")
	}

	// Case: StopFineTunedInstance Fails with no name
	resp, err = req.StopFineTunedInstance(context.TODO(), "")
	if !reflect.DeepEqual(resp, FineTuningResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, FineTuningResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "no name provided")
	}

	// Case: StopFineTunedInstance Fails with invalid HTTP request
	resp, err = req.StopFineTunedInstance(context.TODO(), "a")
	if !reflect.DeepEqual(resp, FineTuningResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, FineTuningResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "!nil")
	}

	// Case: StopFineTunedInstance Fails with HTTP 200 and malformed body
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))

	req.BaseURL = ts.URL

	resp, err = req.StopFineTunedInstance(context.TODO(), "a")
	if !reflect.DeepEqual(resp, FineTuningResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, FineTuningResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "!nil")
	}

	ts.Close()

	// Case: StopFineTunedInstance Fails with HTTP 400
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		fmt.Fprintln(w, nil)
	}))

	req.BaseURL = ts.URL

	resp, err = req.StopFineTunedInstance(context.TODO(), "a")
	if !reflect.DeepEqual(resp, FineTuningResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, FineTuningResponse{})
	}
	if err == nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "!nil")
	}

	ts.Close()

	// Case: StopFineTunedInstance Succeeds with HTTP 200
	req.Debug = true
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, _ := json.Marshal(EmbeddingsResponse{})
		fmt.Fprintln(w, bytes.NewBuffer(resp))
	}))

	req.BaseURL = ts.URL

	resp, err = req.StopFineTunedInstance(context.TODO(), "a")
	if !reflect.DeepEqual(resp, FineTuningResponse{}) {
		t.Errorf("Result was incorrect, got: %v, want: %v.", resp, FineTuningResponse{})
	}
	if err != nil {
		t.Errorf("Error was incorrect, got: %s, want: %s.", err, "!nil")
	}

	ts.Close()
}
