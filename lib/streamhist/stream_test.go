package streamhist

import (
	"math"
	"testing"
)

func TestStreamSummary(t *testing.T) {
	testSz := 1000
	testErr := 0.001

	// test increasing numbers
	s := NewStreamSummaryStructure(testErr)
	for i := 1; i < testSz; i++ {
		s.Update(float64(20 * i))
	}

	println("validating")
	// get each increasing number
	h := s.Histogram()
	errSpan := float64(testSz) * testErr * 20
	for i := 0; i < testSz; i++ {
		quant := h.Quantile(i + 1)
		expect := float64(20 * (i + 1))

		if math.Abs(quant-expect) > errSpan {
			t.Errorf("rank=%d expected=%f received=%2.1f, allowable error=%2.3f [rankspan=%d]", i+1, expect, quant, testErr, int(errSpan))
		}
	}
}

func BenchmarkStreamSummary(b *testing.B) {
	testSz := b.N
	testErr := 0.001
	s := NewStreamSummaryStructure(testErr)

	if testSz < 1000 {
		return
	}

	b.ResetTimer()
	for i := 0; i < testSz; i++ {
		s.Update(float64(i))
	}
}
