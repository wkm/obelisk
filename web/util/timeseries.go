package util

import (
	"log"
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
func DownSample(resolution uint, data []*DataPoint) []SampledDataPoint {
	samples := make([]SampledDataPoint, resolution)
	if len(data) < 1 {
		return samples
	}

	startTime := (*data[0]).Time()
	endTime := (*data[len(data)-1]).Time()

	delTime := (endTime - startTime) / uint64(resolution)

	dataCursor := 0
	for bucketCursor := uint(0); bucketCursor < resolution; bucketCursor++ {

		bucketStart := startTime + uint64(bucketCursor)*delTime
		bucketEnd := bucketStart + delTime

		sampleCount := 1
		var mean, stddev float64 = 0, 0
		for dataCursor < len(samples) && (*data[dataCursor]).Time() < bucketEnd {
			sample := (*data[dataCursor]).Value()
			tmpMean := mean
			mean = mean + (sample-tmpMean)/float64(sampleCount)
			stddev = stddev + (sample-tmpMean)*(sample-tmpMean)
			sampleCount++
			dataCursor++
		}
		if sampleCount > 1 {
			stddev = math.Sqrt(stddev / (float64(sampleCount) - 1))
		}

		log.Printf("%d mean: %f stddev: %f", bucketStart, mean, stddev)

		samples[bucketCursor] = SampledDataPoint{bucketStart, mean, stddev}
	}

	return samples
}
