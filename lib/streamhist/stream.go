package streamhist

type StreamSummaryStructure struct {
	Err   float64
	Count int

	baseSize  int
	summaries []*[]elem
	head      *SummaryStructure
}

func NewStreamSummaryStructure(err float64) *StreamSummaryStructure {
	s := StreamSummaryStructure{}
	s.Err = err
	s.baseSize = 10
	s.summaries = make([]*[]elem, 0, 20)
	s.head = NewSummaryStructure(s.baseSize, err)
	return &s
}

func (s *StreamSummaryStructure) Update(v float64) {
	s.head.Update(v)

	// if the head is now full
	if s.head.Count == s.head.Width {
		// println("stream compress")
		summary := compress(s.head.Histogram().S, s.head.Width, s.head.Err/2)
		s.summaries = append(s.summaries, &summary)

		s.head = NewSummaryStructure(s.baseSize*(1<<uint(len(s.summaries))), s.Err)
	}
}

func (s *StreamSummaryStructure) Histogram() *Histogram {
	summary := compress(s.head.Histogram().S, s.head.Width, s.head.Err/2)
	for _, sum := range s.summaries {
		println("merging: ", len(*sum), (*sum)[0].String(), (*sum)[len(*sum)-1].String())
		summary = merge(summary, *sum)
		println("  -- ", len(summary))
	}

	h := Histogram{}
	h.S = summary
	maxrank := summary[len(summary)-1].rmax
	println("initial maxrank = ", maxrank)
	h.Rank = maxrank + int(s.Err*float64(maxrank))
	println("  with error = ", maxrank)
	h.Err = s.Err

	println("Histogram has max rank", h.Rank, "and", len(h.S), "datapoints")

	// for _, p := range h.S {
	// 	println("  ", p.String())
	// }

	return &h
}
