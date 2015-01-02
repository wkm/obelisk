package util

import (
	"math"
)

// StdDevStream gives the standard deviation across data points.
type StdDevStream struct {
	M0, M1, M2 float64
}

// Sample updates the stream with a given value.
func (s *StdDevStream) Sample(value float64) {
	s.M0++
	s.M1 += value
	s.M2 += math.Pow(value, 2)
}

// Count gives the number of samples seen.
func (s *StdDevStream) Count() uint {
	return uint(s.M0)
}

// Sum gives the sum of all samples seen.
func (s *StdDevStream) Sum() float64 {
	return s.M1
}

// Mean gives the mean of all samples seen.
func (s *StdDevStream) Mean() float64 {
	// no elements
	if s.M0 == 0 {
		return 0
	}

	return s.M1 / s.M0
}

// SampleVariance gives the sample variance of all samples seen.
func (s *StdDevStream) SampleVariance() float64 {
	// no elements
	if s.M0 == 0 {
		return 0
	}

	return (s.M0*s.M2 - s.M1*s.M1) / (s.M0 * (s.M0 - 1))
}

// SampleStdDev gives the standard deviation of all samples seen.
func (s *StdDevStream) SampleStdDev() float64 {
	return math.Sqrt(s.SampleVariance())
}
