package util

import (
	"math"
	"testing"
)

func eq(a, b float64) bool {
	del := math.Abs(a - b)
	return del < 0.0001
}

func TestStdDevStream(t *testing.T) {
	t.Parallel()

	s := StdDevStream{}

	if s.Mean() != 0 {
		t.Errorf("expected mean=0, got %f", s.Mean())
	}
	if s.SampleVariance() != 0 {
		t.Errorf("expected var=0, got %f", s.SampleVariance())
	}
	if s.SampleStdDev() != 0 {
		t.Errorf("expected stddev=0, got %f", s.SampleStdDev())
	}

	s.Sample(2)
	s.Sample(4)
	s.Sample(4)
	s.Sample(4)
	s.Sample(5)
	s.Sample(5)
	s.Sample(7)
	s.Sample(9)

	if s.Count() != 8 {
		t.Errorf("expected count=8, got %d", s.Count())
	}
	if s.Sum() != 40 {
		t.Errorf("expected sum=40, got %f", s.Sum())
	}
	if s.Mean() != 5 {
		t.Errorf("expected mean=5, got %f", s.Mean())
	}
	if !eq(s.SampleVariance(), 4.571429) {
		t.Errorf("expected var=4.571429, got %f", s.SampleVariance())
	}
	if !eq(s.SampleStdDev(), 2.138090) {
		t.Errorf("expected stddev=2.138090, got %f", s.SampleStdDev())
	}
}
