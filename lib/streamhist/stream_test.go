package streamhist

import (
	"math"
	"testing"
)

func TestStreamSummary(t *testing.T) {
	testSz := 1000000
	testErr := 0.001

	// test increasing numbers
	s := NewStreamSummaryStructure(testErr)
	for i := 0; i < testSz; i++ {
		s.Update(float64(20 * i))
	}

	println("validating")
	// get each increasing number
	h := s.Histogram()
	errSpan := float64(testSz) * testErr * 20
	for i := 0; i < testSz; i++ {
		quant := h.Quantile(i + 1)
		expect := float64(20 * i)

		if math.Abs(quant-expect) > errSpan {
			t.Errorf("expected %f; received %2.1f, exceeding allowable error of %2.3f [%d]", expect, quant, testErr, int(errSpan))
			return
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
