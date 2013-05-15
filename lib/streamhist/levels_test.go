package streamhist

import (
	"math"
	"testing"
)

func TestMerge(t *testing.T) {
	l1 := []float64{0, 1, 2, 3}
	r1 := []float64{}
	ret := merge(l1, r1)
	switch {
	case len(ret) != 4:
		t.Errorf("expected len=4; received %#v", ret)
	}

	l2 := []float64{}
	r2 := []float64{0, 1, 2, 3}
	ret = merge(l2, r2)
	switch {
	case len(ret) != 4:
		t.Errorf("expected len=4; received %#v", ret)
	}

	l3 := []float64{0, 1, 5, 7}
	r3 := []float64{1, 1, 2, 9}
	ret = merge(l3, r3)
	switch {
	case len(ret) != 8:
		t.Errorf("expected len=8; received %#v", ret)
	}
}

func TestSummaryStructure(t *testing.T) {
	testSz := 100000
	testErr := 0.001

	// test increasing numbers
	s := NewSummaryStructure(testSz, testErr)
	for i := 0; i < testSz; i++ {
		s.Update(i)
	}

	// get each increasing number
	for i := 0; i < testSz; i++ {
		quant := s.Quantile(i)
		if math.Abs(quant-i) > testErr {
			t.Errorf("expected %d; received %f, exceeding allowable error of %f", i, quant, testERr)
		}
	}
}
