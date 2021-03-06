package streamhist

import (
	"fmt"
)

// Histogram is a frozen, queryable summary structure
type Histogram struct {
	Err  float64 // allowed error
	S    []elem  // a summary sketch
	Rank int     // the maximum rank encoded by the histogram
}

// Quantile extracts a quantile from the histogram
func (h *Histogram) Quantile(rank int) float64 {
	if rank <= 0 || rank > h.Rank {
		panic(fmt.Sprintf("given rank %d is outside of range %d", rank, h.Rank))
	}

	return quantile(h.S, rank, h.Err).val
}
