package main

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

var sampleScored = []scoredPair{
	{pair{"a", "b", false}, 0.95},
	{pair{"c", "d", true}, 0.90},
	{pair{"e", "f", false}, 0.85},
	{pair{"g", "h", true}, 0.80},
}

func TestEvaluateThreshold(t *testing.T) {
	tests := []struct {
		threshold   float64
		sameCaught  int
		trapsCaught int
	}{
		{0.99, 0, 0},
		{0.95, 0, 1},
		{0.90, 1, 1},
		{0.80, 2, 2},
	}

	for _, tt := range tests {
		got := evaluateThreshold(sampleScored, tt.threshold)
		if got.sameCaught != tt.sameCaught || got.trapsCaught != tt.trapsCaught {
			t.Errorf("threshold %.2f: got same %d, traps %d; want same %d, traps %d",
				tt.threshold, got.sameCaught, got.trapsCaught, tt.sameCaught, tt.trapsCaught)
		}
		if got.sameTotal != 2 || got.trapsTotal != 2 {
			t.Errorf("threshold %.2f: got totals %d/%d, want 2/2", tt.threshold, got.sameTotal, got.trapsTotal)
		}
	}
}

func TestSeparation(t *testing.T) {
	lowestSame, highestTrap := separation(sampleScored)
	if lowestSame != 0.80 || highestTrap != 0.95 {
		t.Errorf("got lowestSame %v, highestTrap %v; want 0.80, 0.95", lowestSame, highestTrap)
	}
}

func TestPrintReportDetectsOverlap(t *testing.T) {
	var out bytes.Buffer
	err := printReport(&out, report{
		model:      "test-model",
		scored:     sampleScored,
		thresholds: []float64{0.9},
		latency:    42 * time.Millisecond,
		samples:    3,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, want := range []string{"test-model", "Ranges overlap", "42ms"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("report does not contain %q:\n%s", want, out.String())
		}
	}
}
