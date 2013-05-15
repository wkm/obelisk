package streamhist

import (
	"math"
)

// the multi-level summary structure described by the paper
type SummaryStructure struct {
	Width uint64
	Err   float64
	S     []*Summary

	blockSize  int // `b` in the paper
	levelCount int // `L` in the paper
}

type Summary struct {
	Min, Max float64
	Data     []float64
}

// compute the size blocks we need to cover a fixed amount of data
// within the given error
func computeBlockSize(dataSize uint64, err float64) int {
	// \lfloor \log(err * dataSize) / err \rfloor
	return int(math.Log2(err*float64(dataSize)) / err)
}

func computeLevels(dataSize uint64, blockSize int) int {
	return int(math.Ceil(math.Log2(float64(dataSize) * float64(blockSize))))
}

func NewSummaryStructure(width uint64, err float64) *SummaryStructure {
	blockSize := computeBlockSize(width, err)
	levelCount := computeLevels(width, blockSize)

	// each level_i in the summary
	summarySplay := make([]*Summary, 0, levelCount)
	for i := 0; i < levelCount; i++ {
		summarySplay[i] = &Summary{math.MaxFloat64, -math.MaxFloat64, make([]float64, 0, blockSize)}
	}

	return &SummaryStructure{width, err, summarySplay, blockSize, levelCount}
}

func (s *SummaryStructure) Update(value float64) {
	s0 := append(*s.S[0], value)
	s.S[0] = &s0

	var sc Summary
	if len(s.S[0].Data) == cap(s.S[0].Data) {
		sc = compress(s.S[0], 1/s.blockSize)
		empty(s.S[0])
	}

	for i := 0; i < s.levelCount; i++ {
		if s.S[i].Size() == 0 {
			s.S[0] = sc
			break
		} else {
			sc = s.compress(s.merge(s.S[i], sc), 1/s.blockSize)
			s.S[i].Empty()
		}
	}
}

// EMPTY() in the paper
func empty(s *Summary) {
	// without reallocating the backing array, reset this slice
	// to length 0
	s.Max = -math.MaxFloat64
	s.Min = math.MaxFloat64
	s.Data = s.Data[:0]
}

// COMPRESS() in the paper; but instead of 1/b we use b directly;
// note that the summary -must- be sorted already
func compress(data []float64, width, totalwidth int) []float64 {
	count := math.Ceil(width/2) + 1
	newdata := make([]float64, count)
	for i := 0; i < count; i++ {
		... need to take into account rmin/rmax
		position := int(i * (2 * totalwidth / width))
		newdata[i] = data[position]
	}
	return newdata
}

// MERGE() in the paper; give a new array that
func merge(left, right []float64) []float64 {
	result := make([]float64, len(left)+len(right))

	i, k := 0, 0
	for {
		// if nothing left on the left; zip through right
		if i == len(left) {
			for k < len(right) {
				result[i+k] = right[k]
				k++
			}
			break
		}

		// if nothing left on the right; zip through left
		if k == len(right) {
			for i < len(left) {
				result[i+k] = left[i]
				i++
			}
			break
		}

		if left[i] < right[k] {
			i++
		}
	}
}
