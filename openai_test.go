package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeItem struct {
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
}

func newTestClient(t *testing.T, handler http.HandlerFunc) *embeddingClient {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client := newEmbeddingClient("test-key", "test-model")
	client.baseURL = server.URL

	return client
}

func fakeEmbeddings(t *testing.T, vectors map[string][]float64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("got Authorization %q", got)
		}

		var req embeddingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decode request: %v", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if req.Model != "test-model" {
			t.Errorf("got model %q, want test-model", req.Model)
		}

		items := make([]fakeItem, 0, len(req.Input))
		for i := len(req.Input) - 1; i >= 0; i-- {
			items = append(items, fakeItem{Index: i, Embedding: vectors[req.Input[i]]})
		}
		json.NewEncoder(w).Encode(map[string][]fakeItem{"data": items})
	}
}

func TestEmbedPreservesInputOrder(t *testing.T) {
	client := newTestClient(t, fakeEmbeddings(t, map[string][]float64{
		"first":  {1, 0},
		"second": {0, 1},
	}))

	got, err := client.embed(context.Background(), []string{"first", "second"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got[0][0] != 1 || got[1][1] != 1 {
		t.Errorf("embeddings are out of order: %v", got)
	}
}

func TestEmbedReturnsAPIError(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"error":"invalid api key"}`, http.StatusUnauthorized)
	})

	_, err := client.embed(context.Background(), []string{"question"})
	if err == nil || !strings.Contains(err.Error(), "401") || !strings.Contains(err.Error(), "invalid api key") {
		t.Errorf("got error %v, want status 401 with API message", err)
	}
}

func TestEmbedRejectsWrongCount(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string][]fakeItem{"data": {{Index: 0, Embedding: []float64{1}}}})
	})

	_, err := client.embed(context.Background(), []string{"first", "second"})
	if err == nil || !strings.Contains(err.Error(), "got 1 embeddings for 2 inputs") {
		t.Errorf("got error %v, want count mismatch", err)
	}
}

func TestScorePairsSortsBySimilarity(t *testing.T) {
	client := newTestClient(t, fakeEmbeddings(t, map[string][]float64{
		"enable":  {1, 0.1},
		"disable": {1, 0.2},
		"reset":   {1, 0},
		"forgot":  {0, 1},
	}))

	scored, err := scorePairs(context.Background(), client, []pair{
		{"reset", "forgot", true},
		{"enable", "disable", false},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if scored[0].first != "enable" || scored[1].first != "reset" {
		t.Errorf("pairs are not sorted by similarity: %+v", scored)
	}
}
