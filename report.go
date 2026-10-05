package main

import (
	"fmt"
	"io"
	"math"
	"text/tabwriter"
	"time"
)

type report struct {
	model      string
	scored     []scoredPair
	thresholds []float64
	latency    time.Duration
	samples    int
}

type thresholdStats struct {
	threshold   float64
	sameCaught  int
	sameTotal   int
	trapsCaught int
	trapsTotal  int
}

func evaluateThreshold(scored []scoredPair, threshold float64) thresholdStats {
	stats := thresholdStats{threshold: threshold}
	for _, p := range scored {
		isHit := p.similarity >= threshold
		if p.sameMeaning {
			stats.sameTotal++
			if isHit {
				stats.sameCaught++
			}
			continue
		}
		stats.trapsTotal++
		if isHit {
			stats.trapsCaught++
		}
	}

	return stats
}

func separation(scored []scoredPair) (lowestSame, highestTrap float64) {
	lowestSame, highestTrap = math.Inf(1), math.Inf(-1)
	for _, p := range scored {
		if p.sameMeaning {
			lowestSame = min(lowestSame, p.similarity)
			continue
		}
		highestTrap = max(highestTrap, p.similarity)
	}

	return lowestSame, highestTrap
}

func label(p pair) string {
	if p.sameMeaning {
		return "SAME"
	}

	return "TRAP"
}

func printReport(w io.Writer, r report) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintf(tw, "Model: %s\n\n", r.model)

	fmt.Fprintln(tw, "SIMILARITY\tTYPE\tFIRST\tSECOND")
	for _, p := range r.scored {
		fmt.Fprintf(tw, "%.3f\t%s\t%s\t%s\n", p.similarity, label(p.pair), p.first, p.second)
	}
	fmt.Fprintln(tw)

	fmt.Fprintln(tw, "THRESHOLD\tSAME CAUGHT (useful hits)\tTRAPS CAUGHT (wrong answers)")
	for _, threshold := range r.thresholds {
		s := evaluateThreshold(r.scored, threshold)
		fmt.Fprintf(tw, "%.2f\t%d/%d\t%d/%d\n", s.threshold, s.sameCaught, s.sameTotal, s.trapsCaught, s.trapsTotal)
	}
	fmt.Fprintln(tw)

	lowestSame, highestTrap := separation(r.scored)
	fmt.Fprintf(tw, "Lowest SAME pair:  %.3f\n", lowestSame)
	fmt.Fprintf(tw, "Highest TRAP pair: %.3f\n", highestTrap)
	if highestTrap >= lowestSame {
		fmt.Fprintln(tw, "Ranges overlap: no threshold separates paraphrases from traps.")
	} else {
		fmt.Fprintf(tw, "Any threshold in (%.3f, %.3f] separates them on this dataset.\n", highestTrap, lowestSame)
	}
	fmt.Fprintf(tw, "Median embedding latency, 1 question, %d samples: %s\n", r.samples, r.latency.Round(time.Millisecond))

	return tw.Flush()
}
