// Command semantic-cache-experiment checks whether a similarity threshold
// can tell paraphrases apart from look-alike questions.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
)

const defaultModel = "text-embedding-3-small"

var thresholds = []float64{0.70, 0.75, 0.80, 0.85, 0.90, 0.95}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	model := flag.String("model", defaultModel, "OpenAI embedding model")
	samples := flag.Int("samples", 5, "number of requests used to measure embedding latency")
	flag.Parse()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return errors.New("OPENAI_API_KEY is not set")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	client := newEmbeddingClient(apiKey, *model)

	scored, err := scorePairs(ctx, client, defaultPairs)
	if err != nil {
		return err
	}

	latency, err := measureLatency(ctx, client, *samples)
	if err != nil {
		return err
	}

	return printReport(os.Stdout, report{
		model:      *model,
		scored:     scored,
		thresholds: thresholds,
		latency:    latency,
		samples:    *samples,
	})
}
