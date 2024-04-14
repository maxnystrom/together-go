package together

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CompletionsRequest struct {
	Model             string   `json:"model"`
	Prompt            string   `json:"prompt"`
	MaxTokens         int32    `json:"max_tokens"`
	Stream            bool     `json:"stream"`
	Stop              []string `json:"stop"`
	Temperature       float64  `json:"temperature"`
	TopP              float64  `json:"top_p"`
	TopK              int32    `json:"top_k"`
	RepetitionPenalty float64  `json:"repetition_penalty"`
	Logprobs          int32    `json:"logprobs"`
	Echo              bool     `json:"echo"`
	N                 int32    `json:"n"`
	SafetyModel       string   `json:"safety_model"`
}

type CompletionsResponse struct {
	Id      string         `json:"id"`
	Choices []ChoiceObject `json:"choices"`
	Usage   UsageObject    `json:"usage"`
	Created int            `json:"created"`
	Model   string         `json:"model"`
	Object  string         `json:"object"`
}

type ChoiceObject struct {
	Text string `json:"text"`
}

// Completions is the endpoint for language, code, and image models on Together AI.
//
// API Reference: https://docs.together.ai/reference/completions
func (api *API) Completions(ctx context.Context, model string, prompt string, maxTokens int32, request CompletionsRequest) (CompletionsResponse, error) {
	if ctx == nil {
		return CompletionsResponse{}, fmt.Errorf("no context provided")
	}
	if model == "" {
		return CompletionsResponse{}, fmt.Errorf("no model provided")
	}
	if prompt == "" {
		return CompletionsResponse{}, fmt.Errorf("no prompt provided")
	}
	if maxTokens < 1 {
		return CompletionsResponse{}, fmt.Errorf("maxTokens must be greater than 0 and less than 2147483647") // Not sure it makes sense to have negative tokens?
	}

	request.Model = model
	request.Prompt = prompt
	request.MaxTokens = maxTokens

	uri := defaultBasePath + Version + "/completions"
	reqBody, err := json.Marshal(request)
	if err != nil {
		return CompletionsResponse{}, err
	}

	res, err := api.request(ctx, "POST", uri, bytes.NewBuffer(reqBody), nil)
	if err != nil {
		return CompletionsResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return CompletionsResponse{}, err
	}
	if res.StatusCode != http.StatusOK {
		return CompletionsResponse{}, fmt.Errorf("HTTP request failed: %s", string(body))
	}

	var completionsResponse CompletionsResponse
	err = json.Unmarshal(body, &completionsResponse)
	if err != nil {
		return CompletionsResponse{}, err
	}

	return completionsResponse, nil
}
