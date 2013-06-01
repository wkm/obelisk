package rinst

import (
	"obelisk/lib/streamhist"
)

var DefaultPercentiles = []float64{
	0, 25, 50, 75, 90, 99, 99.9, 99.99, 100,
}

// ... default percentiles: [0, 25, 50, 75, 90, 99, 99.9, 99.99, 100]
type FloatStream struct {
	desc, unit  string
	precentiles []float64 // which percentiles to record

	s *streamhist.StreamSummaryStructure
}

func NewFloatStream(desc, unit string, err float64) *FloatStream {
	fs := FloatStream{}
	fs.desc = desc
	fs.unit = unit
	fs.precentiles = DefaultPercentiles
	fs.s = streamhist.NewStreamSummaryStructure(err)
	return fs
}

func (s *FloatStream) Record(value float64) {
	s.s.Update(value)
}

func (s *FloatStream) Measure(n string, b MeasurementBuffer) {
	histo := s.s.Histogram()
	now := uint64(time.Now.Unix())
	for _, p := range s.precentiles {
		val := histo.Quantile(histo.Rank * p / 100)
		b <- Measurement{fmt.Sprintf("%s_%f", name, p), fmt.Sprintf("%f", val)}
	}
}

func (s *FloatStream) Schema(n string, b SchemaBuffer) {
	b <- Schema{name, TypeFloatStream, s.unit, s.desc}
}
