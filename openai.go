package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultBaseURL    = "https://api.openai.com/v1"
	requestTimeout    = 30 * time.Second
	maxErrorBodyBytes = 4 << 10
)

type embeddingClient struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

func newEmbeddingClient(apiKey, model string) *embeddingClient {
	return &embeddingClient{
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
		model:      model,
		httpClient: &http.Client{Timeout: requestTimeout},
	}
}

type embeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type embeddingResponse struct {
	Data []struct {
		Index     int       `json:"index"`
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
}

func (c *embeddingClient) embed(ctx context.Context, inputs []string) ([][]float64, error) {
	body, err := json.Marshal(embeddingRequest{Model: c.model, Input: inputs})
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBodyBytes))
		return nil, fmt.Errorf("openai: status %d: %s", resp.StatusCode, bytes.TrimSpace(msg))
	}

	var parsed embeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if len(parsed.Data) != len(inputs) {
		return nil, fmt.Errorf("openai: got %d embeddings for %d inputs", len(parsed.Data), len(inputs))
	}

	embeddings := make([][]float64, len(inputs))
	for _, item := range parsed.Data {
		if item.Index < 0 || item.Index >= len(inputs) || embeddings[item.Index] != nil {
			return nil, fmt.Errorf("openai: unexpected embedding index %d", item.Index)
		}
		embeddings[item.Index] = item.Embedding
	}

	return embeddings, nil
}
