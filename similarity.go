package main

import (
	"errors"
	"fmt"
	"math"
)

var (
	errDimensionMismatch = errors.New("vectors have different dimensions")
	errZeroVector        = errors.New("cosine similarity is undefined for a zero vector")
)

func cosineSimilarity(a, b []float64) (float64, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("%w: %d vs %d", errDimensionMismatch, len(a), len(b))
	}

	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0, errZeroVector
	}

	return dot / (math.Sqrt(normA) * math.Sqrt(normB)), nil
}
