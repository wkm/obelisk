package rinst

// an instrument can be measured
type Instrument interface {
	Measure() string
}
