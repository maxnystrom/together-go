package together

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ChatCompletionsRequest struct {
	Model             string               `json:"model"`
	Messages          []Message            `json:"messages"`
	Stream            bool                 `json:"stream"`
	MaxTokens         int32                `json:"max_tokens"`
	Stop              []string             `json:"stop"`
	Temperature       float64              `json:"temperature"`
	TopP              float64              `json:"top_p"`
	TopK              int32                `json:"top_k"`
	RepetitionPenalty float64              `json:"repetition_penalty"`
	Logprobs          int32                `json:"logprobs"`
	Echo              bool                 `json:"echo"`
	N                 int32                `json:"n"`
	SafetyModel       string               `json:"safety_model"`
	ResponseFormat    ResponseFormatObject `json:"response_format"`
	Tools             []Tool               `json:"tools"`
	ToolChoice        []ToolChoiceObject   `json:"tool_choice"`
	FrequencyPenalty  float64              `json:"frequency_penalty"`
	PresencePenalty   float64              `json:"presence_penalty"`
	MinP              float64              `json:"min_p"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseFormatObject struct {
	Type   string          `json:"type"`
	Schema json.RawMessage `json:"schema"` // maybe should be *json.RawMessage ?
}

type Tool struct {
	Type     string         `json:"type"`
	Function FunctionObject `json:"function"`
}

type FunctionObject struct {
	Description string          `json:"description"`
	Name        string          `json:"name"`
	Parameters  json.RawMessage `json:"parameters"` // maybe should be *json.RawMessage ?
}

type ToolChoiceObject struct {
	Type     string         `json:"type"`
	Function FunctionObject `json:"function"`
}

type ChatCompletionsResponse struct {
	Id      string      `json:"id"`
	Choices []Message   `json:"choices"`
	Usage   UsageObject `json:"usage"`
	Created int         `json:"created"`
	Model   string      `json:"model"`
	Object  string      `json:"object"`
}

type UsageObject struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Chat Completions is the endpoint for chat and moderation models on Together AI.
//
// API Reference: https://docs.together.ai/reference/chat-completions
func (api *API) ChatCompletions(ctx context.Context, model string, messages []Message, request ChatCompletionsRequest) (ChatCompletionsResponse, error) {
	if len(messages) == 0 || messages == nil {
		return ChatCompletionsResponse{}, fmt.Errorf("no messages provided")
	}
	if model == "" {
		return ChatCompletionsResponse{}, fmt.Errorf("no model provided")
	}
	if ctx == nil {
		return ChatCompletionsResponse{}, fmt.Errorf("no context provided")
	}

	request.Messages = messages
	request.Model = model

	uri := defaultBasePath + Version + "/chat/completions"

	res, err := api.request(ctx, "POST", uri, nil, nil)
	if err != nil {
		return ChatCompletionsResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ChatCompletionsResponse{}, err
	}
	if res.StatusCode != http.StatusOK {
		return ChatCompletionsResponse{}, fmt.Errorf("HTTP request failed: %s", string(body))
	}

	var chatCompletionsResponse ChatCompletionsResponse
	err = json.Unmarshal(body, &chatCompletionsResponse)
	if err != nil {
		return ChatCompletionsResponse{}, err
	}

	return chatCompletionsResponse, nil
}
