package main

import (
	"errors"
	"math"
	"testing"
)

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name string
		a, b []float64
		want float64
	}{
		{"identical", []float64{1, 2, 3}, []float64{1, 2, 3}, 1},
		{"scaled", []float64{1, 2, 3}, []float64{2, 4, 6}, 1},
		{"orthogonal", []float64{1, 0}, []float64{0, 1}, 0},
		{"opposite", []float64{1, 2}, []float64{-1, -2}, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cosineSimilarity(tt.a, tt.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCosineSimilarityErrors(t *testing.T) {
	tests := []struct {
		name    string
		a, b    []float64
		wantErr error
	}{
		{"dimension mismatch", []float64{1, 2}, []float64{1, 2, 3}, errDimensionMismatch},
		{"zero vector", []float64{0, 0}, []float64{1, 2}, errZeroVector},
		{"empty vectors", nil, nil, errZeroVector},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cosineSimilarity(tt.a, tt.b)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got error %v, want %v", err, tt.wantErr)
			}
		})
	}
}
