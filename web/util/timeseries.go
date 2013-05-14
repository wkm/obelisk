package util

import (
	"math"
)

type DataPoint interface {
	Time() uint64
	Value() float64
}

type DSPoint struct {
	T uint64
	V float64
}

func (d DSPoint) Time() uint64   { return d.T }
func (d DSPoint) Value() float64 { return d.V }

type SampledDataPoint struct {
	Time     uint64
	Avg, Err float64
}

// down sample a sorted time series into a target resolution
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
