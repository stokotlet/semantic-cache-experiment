package main

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"
	"time"
)

const latencyProbe = "How do I reset my password?"

type scoredPair struct {
	pair
	similarity float64
}

func scorePairs(ctx context.Context, client *embeddingClient, pairs []pair) ([]scoredPair, error) {
	inputs := make([]string, 0, len(pairs)*2)
	for _, p := range pairs {
		inputs = append(inputs, p.first, p.second)
	}

	embeddings, err := client.embed(ctx, inputs)
	if err != nil {
		return nil, fmt.Errorf("embed pairs: %w", err)
	}

	scored := make([]scoredPair, len(pairs))
	for i, p := range pairs {
		similarity, err := cosineSimilarity(embeddings[2*i], embeddings[2*i+1])
		if err != nil {
			return nil, fmt.Errorf("score %q vs %q: %w", p.first, p.second, err)
		}
		scored[i] = scoredPair{pair: p, similarity: similarity}
	}

	slices.SortFunc(scored, func(a, b scoredPair) int {
		return cmp.Compare(b.similarity, a.similarity)
	})

	return scored, nil
}

func measureLatency(ctx context.Context, client *embeddingClient, samples int) (time.Duration, error) {
	if samples < 1 {
		return 0, errors.New("samples must be at least 1")
	}

	// warm-up: the first request includes the TLS handshake
	if _, err := client.embed(ctx, []string{latencyProbe}); err != nil {
		return 0, fmt.Errorf("warm up: %w", err)
	}

	durations := make([]time.Duration, 0, samples)
	for range samples {
		start := time.Now()
		if _, err := client.embed(ctx, []string{latencyProbe}); err != nil {
			return 0, fmt.Errorf("measure latency: %w", err)
		}
		durations = append(durations, time.Since(start))
	}

	slices.Sort(durations)

	return durations[len(durations)/2], nil
}
