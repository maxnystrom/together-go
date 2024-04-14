package together

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type EmbeddingsRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type EmbeddingsResponse struct {
	Status     string   `json:"status"`
	Prompt     []string `json:"prompt"`
	Model      string   `json:"model"`
	ModelOwner string   `json:"model_owner"`
	Tags       Tags     `json:"tags"`
	NumReturns int      `json:"num_returns"`
	Args       Args     `json:"args"`
	Subjobs    []Subjob `json:"subjobs"`
	Output     Output   `json:"output"`
}

// Embeddings is the endpoint for embedding models on Together AI.
//
// API Reference: https://docs.together.ai/reference/embeddings
func (api *API) Embeddings(ctx context.Context, model string, input string, request EmbeddingsRequest) (EmbeddingsResponse, error) {
	if ctx == nil {
		return EmbeddingsResponse{}, fmt.Errorf("no context provided")
	}
	if model == "" {
		return EmbeddingsResponse{}, fmt.Errorf("no model provided")
	}
	if input == "" {
		return EmbeddingsResponse{}, fmt.Errorf("no input provided")
	}

	request.Model = model
	request.Input = input

	uri := defaultBasePath + Version + "/embeddings"
	reqBody, err := json.Marshal(request)
	if err != nil {
		return EmbeddingsResponse{}, err
	}

	res, err := api.request(ctx, "POST", uri, bytes.NewBuffer(reqBody), nil)
	if err != nil {
		return EmbeddingsResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return EmbeddingsResponse{}, err
	}
	if res.StatusCode != http.StatusOK {
		return EmbeddingsResponse{}, fmt.Errorf("HTTP request failed: %s", string(body))
	}

	var embeddingsResponse EmbeddingsResponse
	err = json.Unmarshal(body, &embeddingsResponse)
	if err != nil {
		return EmbeddingsResponse{}, err
	}

	return embeddingsResponse, nil
}
