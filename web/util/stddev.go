package util

import (
	"math"
)

type StdDevStream struct {
	M0, M1, M2 float64
}

func (s *StdDevStream) Sample(value float64) {
	s.M0++
	s.M1 += value
	s.M2 += math.Pow(value, 2)
}

func (s *StdDevStream) Count() uint {
	return uint(s.M0)
}

func (s *StdDevStream) Sum() float64 {
	return s.M1
}

func (s *StdDevStream) Mean() float64 {
	// no elements
	if s.M0 == 0 {
		return 0
	}

	return s.M1 / s.M0
}

// give the sample variance computed from the stream
func (s *StdDevStream) SampleVariance() float64 {
	// no elements
	if s.M0 == 0 {
		return 0
	}

	return (s.M0*s.M2 - s.M1*s.M1) / (s.M0 * (s.M0 - 1))
}

func (s *StdDevStream) SampleStdDev() float64 {
	return math.Sqrt(s.SampleVariance())
}
