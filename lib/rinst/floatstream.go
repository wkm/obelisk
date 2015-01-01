package rinst

import (
	"fmt"
	"time"

	"github.com/wkm/obelisk/lib/streamhist"
)

// DefaultPercentiles contains the percentiles which are measured by default
// in this stream of floats.
var DefaultPercentiles = []float64{
	0, 25, 50, 75, 90, 99, 99.9, 99.99, 100,
}

// FloatStream is an instrument for efficiently measuring percentiles from a
// stream of float values.
type FloatStream struct {
	desc, unit  string
	precentiles []float64 // which percentiles to record

	s *streamhist.StreamSummaryStructure
}

// NewFloatStream allocates an instrument for computing the histogram of a
// stream of floats.
func NewFloatStream(desc, unit string, err float64) *FloatStream {
	fs := FloatStream{}
	fs.desc = desc
	fs.unit = unit
	fs.precentiles = DefaultPercentiles
	fs.s = streamhist.NewStreamSummaryStructure(err)
	return &fs
}

// Record saves a value into a stream of floats.
func (s *FloatStream) Record(value float64) {
	s.s.Update(value)
}

// Measure writes the percentiles of this float stream into
func (s *FloatStream) Measure(n string, r MeasurementReceiver) {
	histo := s.s.Histogram()
	now := time.Now().Unix()
	for _, p := range s.precentiles {
		rank := float64(histo.Rank) * (p / 100.0)
		val := histo.Quantile(int(rank))
		r.WriteFloat(fmt.Sprintf("%s_%f", n, p), now, val)
	}
}

// Schema writes the recorded percentiles as individual metrics into
// the given receiver.
func (s *FloatStream) Schema(n string, r SchemaReceiver) {
	r.WriteSchema(n, TypeFloatStream, s.unit, s.desc)
}
