package util

import (
	"math"
)

// DataPoint represents a type which has a timestamp and a value.
type DataPoint interface {
	Time() uint64
	Value() float64
}

// DSPoint is a trivial implementation of a DataPoint.
type DSPoint struct {
	T uint64
	V float64
}

// Time gives the time associated with a data point.
func (d DSPoint) Time() uint64 { return d.T }

// Value gives value from a data point.
func (d DSPoint) Value() float64 { return d.V }

// SampledDataPoint is a datapoint with an error value as well.
type SampledDataPoint struct {
	Time     uint64
	Avg, Err float64
}

// DownSample a sorted time series into a target resolution, keeping track of error.
func DownSample(start, stop uint64, resolution uint, data []*DataPoint) []SampledDataPoint {
	samples := make([]SampledDataPoint, resolution)
	if len(data) < 1 {
		return samples
	}

	delTime := (stop - start) / uint64(resolution)
	dataCursor := 0
	for bucketCursor := uint(0); bucketCursor < resolution; bucketCursor++ {
		bucketStart := start + uint64(bucketCursor)*delTime
		bucketEnd := bucketStart + delTime

		bucket := StdDevStream{}
		for dataCursor < len(data) && (*data[dataCursor]).Time() < bucketEnd {
			bucket.Sample((*data[dataCursor]).Value())
			dataCursor++
		}

		if bucket.Count() < 1 {
			samples[bucketCursor] = SampledDataPoint{bucketStart, math.NaN(), 0}
		} else {
			samples[bucketCursor] = SampledDataPoint{
				bucketStart,
				bucket.Mean(),
				bucket.SampleStdDev(),
			}
		}
	}

	return samples
}
