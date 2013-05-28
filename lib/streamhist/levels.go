package streamhist

import (
	"fmt"
	"math"
	"sort"
)

// FUTURE
// - s0 does not need rank information

// the multi-level summary structure described by the paper
type SummaryStructure struct {
	Width int
	Err   float64
	S     [][]elem

	blockSize  int // `b` in the paper
	levelCount int // `L` in the paper
}

// a frozen summary structure
type Histogram struct {
	Err  float64
	S    []elem
	Rank int
}

// elem roughly corresponds to a bin in a PDF summary
type elem struct {
	val        float64 // the height of the bin
	rmin, rmax int     // roughly, the boundaries of the histogram bin
}

func (e elem) String() string {
	return fmt.Sprintf("e{%2.1f, [%d,%d]}", e.val, e.rmin, e.rmax)
}

type elemSorter struct {
	e []elem
}

func (e *elemSorter) Len() int {
	return len(e.e)
}
func (e *elemSorter) Swap(i, j int) {
	e.e[i], e.e[j] = e.e[j], e.e[i]
}
func (e *elemSorter) Less(i, j int) bool {
	return e.e[i].val < e.e[j].val
}

func sortElem(e []elem) {
	sort.Sort(&elemSorter{e})
}

// compute the size blocks we need to cover a fixed amount of data
// with the given error
func computeBlockSize(dataSize int, err float64) int {
	if float64(dataSize)*err < 1 {
		return int(dataSize)
	}

	// \lfloor \log(err * dataSize) / err \rfloor
	sz := int(math.Log2(err*float64(dataSize)) / err)
	return sz
}

func computeLevels(dataSize int, blockSize int) int {
	if float64(dataSize)/float64(blockSize) <= 1 {
		return 1
	}

	levels := int(math.Ceil(math.Log2(float64(dataSize) / float64(blockSize))))
	return levels
}

func NewSummaryStructure(width int, err float64) *SummaryStructure {
	blockSize := computeBlockSize(width, err)
	levelCount := computeLevels(width, blockSize) + 1 // we incrememt by one to include s0

	// each level_i in the summary
	summarySplay := make([][]elem, levelCount)
	for i := 0; i < levelCount; i++ {
		summarySplay[i] = make([]elem, 0, blockSize)
	}

	return &SummaryStructure{width, err, summarySplay, blockSize, levelCount}
}

func (s *SummaryStructure) Update(value float64) {
	var sc []elem
	if len(s.S[0]) == cap(s.S[0]) {
		sortElem(s.S[0])

		// we need to assign ranks now
		for i, _ := range s.S[0] {
			s.S[0][i].rmin = i + 1
			s.S[0][i].rmax = i + 1
		}

		sc = compress(s.S[0], s.blockSize, s.Err)
		empty(&s.S[0])

		for i := 1; i < s.levelCount; i++ {
			if len(s.S[i]) == 0 {
				s.S[i] = sc
				break
			} else {
				sc = compress(merge(s.S[i], sc), s.blockSize, s.Err)
				empty(&s.S[i])
			}
		}
	}

	s0 := append(s.S[0], elem{value, 0, 0})
	s.S[0] = s0
}

// EMPTY() in the paper
func empty(e *[]elem) {
	// without reallocating the backing array, reset this slice
	// to length 0
	*e = (*e)[:0]
}

// COMPRESS() in the paper; but instead of 1/b we use b directly;
// note that the summary is assumed sorted
func compress(data []elem, width int, err float64) []elem {
	totalwidth := data[len(data)-1].rmax
	count := int(math.Ceil(float64(width)/2) + 1)
	newdata := make([]elem, count)

	i := 0
	for i < count {
		rank := int(i * (2.0 * totalwidth / width))
		if rank < 1 {
			newdata[i] = quantile(data, 1, err)
			i++
			continue
		}
		if rank > totalwidth {
			newdata[i] = quantile(data, totalwidth, err)
			break
		}

		e := quantile(data, rank, err)

		newdata[i] = e
		i++
	}

	// a little ugly...
	return newdata
}

// quantile() in the paper; give the elem for the given rank
func quantile(data []elem, rank int, err float64) elem {
	if len(data) == 0 {
		return elem{}
	}

	// pick where to start
	lo := 0
	hi := len(data)

	rankErr := int(err * float64(len(data)))
	rankLo := rank - rankErr
	rankHi := rank + rankErr

	bestElem := elem{}
	bestDist := math.MaxInt64

	// FIXME
	// we can skip looking around completly because rank-width is
	// fixed within a range
	for {
		cursor := lo + (hi-lo)/2
		if lo == hi {
			return bestElem
		}

		if cursor < 0 {
			panic("negative cursor")
		}
		if cursor >= len(data) {
			panic("out of bounds cursor")
		}

		elem := data[cursor]
		dist := rank - elem.rmin + (elem.rmax-elem.rmin)/2

		if lo == len(data) {
			panic("couldn't find entry")
		}

		if dist > 0 {
			lo = cursor + 1
		} else if dist <= 0 {
			hi = cursor
		}

		if data[cursor].rmin >= rankLo && data[cursor].rmax <= rankHi {
			if dist < 0 {
				dist = -dist
			}

			if dist < bestDist {
				bestElem = elem
				bestDist = dist
			}
		}
	}

	panic("quantile() panic")
}

// merge two summaries, compressing their
func merge(left, right []elem) []elem {
	result := make([]elem, len(left)+len(right))
	var lastLeft, lastRight elem

	i, k := 0, 0
	for {
		// if nothing left on the left; zip through right
		if i == len(left) {
			for k < len(right) {
				result[i+k] = right[k]
				result[i+k].rmin += lastLeft.rmax
				result[i+k].rmax += lastLeft.rmax
				k++
			}
			break
		}

		// if nothing left on the right; zip through left
		if k == len(right) {
			for i < len(left) {
				result[i+k] = left[i]
				result[i+k].rmin += lastRight.rmax
				result[i+k].rmax += lastRight.rmax
				i++
			}
			break
		}

		if left[i].val < right[k].val {
			lastLeft = left[i]
			result[i+k] = left[i]
			result[i+k].rmin += lastRight.rmax
			result[i+k].rmax += lastRight.rmax
			i++
		} else {
			lastRight = right[k]
			result[i+k] = right[k]
			result[i+k].rmin += lastLeft.rmax
			result[i+k].rmax += lastLeft.rmax
			k++
		}

		// if applicable, adjust the previous element
		// if i+k > 2 {
		// 	result[i+k-2].rmax = result[i+k-1].rmax - 1
		// }
	}

	return result
}

// extract a histogram from the current state of the summary structure
func (s *SummaryStructure) Histogram() *Histogram {
	h := Histogram{}
	h.Err = s.Err
	sortElem(s.S[0])
	for i := range s.S[0] {
		s.S[0][i].rmin = i + 1
		s.S[0][i].rmax = i + 1
	}

	for _, summary := range s.S {
		h.S = merge(h.S, summary)
	}

	h.Rank = h.S[len(h.S)-1].rmax + int(s.Err*float64(len(h.S)))
	return &h
}

// extract a quantile from the histogram
func (h *Histogram) Quantile(rank int) float64 {
	if rank <= 0 || rank > h.Rank {
		panic("given rank is out of range")
	}

	return quantile(h.S, rank, h.Err).val
}
