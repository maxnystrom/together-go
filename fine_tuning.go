package together

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type FineTuningResponse struct {
	Status     string   `json:"status"`
	Prompt     []string `json:"prompt"`
	Model      string   `json:"model"`
	ModelOwner string   `json:"model_owner"`
	Tags       Tags     `json:"tags"` // API Documentation states this is an object, but does not define it or provide an example.
	NumReturns int      `json:"num_returns"`
	Args       Args     `json:"args"`
	Subjobs    []Subjob `json:"subjobs"` // API Documentation states this is an array, but does not define it or provide an example.
	Output     Output   `json:"output"`
}

type Args struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Temperature float64 `json:"temperature"`
	TopP        float64 `json:"top_p"`
	TopK        int     `json:"top_k"`
	MaxTokens   int     `json:"max_tokens"`
}

type Output struct {
	Choices        []Choice `json:"choices"`
	RawComputeTime float64  `json:"raw_compute_time"`
	ResultType     string   `json:"result_type"`
}

type Choice struct {
	FinishReason string `json:"finish_reason"`
	Index        int    `json:"index"`
	Text         string `json:"text"`
}

type Tags struct{}   // TODO: define Tags.
type Subjob struct{} // TODO: define Subjob.

// List Running Instances is the endpoint for listing running fine-tuning instances.
//
// API Reference: https://docs.together.ai/reference/instances
func (api *API) ListRunningInstances(ctx context.Context) (FineTuningResponse, error) {
	if ctx == nil {
		return FineTuningResponse{}, fmt.Errorf("no context provided")
	}

	uri := defaultBasePath + "instances"

	res, err := api.request(ctx, "GET", uri, nil, nil)
	if err != nil {
		return FineTuningResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return FineTuningResponse{}, err
	}
	if res.StatusCode != http.StatusOK {
		return FineTuningResponse{}, fmt.Errorf("HTTP request failed: %s", string(body))
	}

	var fineTuningResponse FineTuningResponse
	err = json.Unmarshal(body, &fineTuningResponse)
	if err != nil {
		return FineTuningResponse{}, err
	}

	return fineTuningResponse, nil
}

// Start Fine-tuned Instance is the endpoint for starting a fine-tuned model.
//
// API Reference: https://docs.together.ai/reference/instances-start
func (api *API) StartFineTunedInstance(ctx context.Context, name string) (FineTuningResponse, error) {
	if ctx == nil {
		return FineTuningResponse{}, fmt.Errorf("no context provided")
	}
	if name == "" {
		return FineTuningResponse{}, fmt.Errorf("no name provided")
	}

	uri := defaultBasePath + "instances/start?model=" + url.QueryEscape(name)

	res, err := api.request(ctx, "POST", uri, nil, nil)
	if err != nil {
		return FineTuningResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return FineTuningResponse{}, err
	}
	if res.StatusCode != http.StatusOK {
		return FineTuningResponse{}, fmt.Errorf("HTTP request failed: %s", string(body))
	}

	var fineTuningResponse FineTuningResponse
	err = json.Unmarshal(body, &fineTuningResponse)
	if err != nil {
		return FineTuningResponse{}, err
	}

	return fineTuningResponse, nil
}

// Stop Fine-tuned Instance is the endpoint for stopping a fine-tuned model.
//
// API Reference: https://docs.together.ai/reference/instances-stop
func (api *API) StopFineTunedInstance(ctx context.Context, name string) (FineTuningResponse, error) {
	if ctx == nil {
		return FineTuningResponse{}, fmt.Errorf("no context provided")
	}
	if name == "" {
		return FineTuningResponse{}, fmt.Errorf("no name provided")
	}

	uri := defaultBasePath + "instances/stop?model=" + url.QueryEscape(name)

	res, err := api.request(ctx, "POST", uri, nil, nil)
	if err != nil {
		return FineTuningResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return FineTuningResponse{}, err
	}
	if res.StatusCode != http.StatusOK {
		return FineTuningResponse{}, fmt.Errorf("HTTP request failed: %s", string(body))
	}

	var fineTuningResponse FineTuningResponse
	err = json.Unmarshal(body, &fineTuningResponse)
	if err != nil {
		return FineTuningResponse{}, err
	}

	return fineTuningResponse, nil
}

// TODO: Looking at the Python CLI, there are clearly undocumented API endpoints for
// Create
// Retrieve
// Monitor
// Status
// Checkpoints
// Download
//
// As these are undocumented, they're likely to change or be removed at any time.
// Therefore adding these endpoints is not a priority.
