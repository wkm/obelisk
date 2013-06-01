package streamhist

import (
	"math"
	"testing"
)

func e(v ...float64) []elem {
	arr := make([]elem, len(v))
	for i, e := range v {
		arr[i] = elem{e, i + 1, i + 1}
	}
	return arr
}

func TestSorter(t *testing.T) {
	e := []elem{
		elem{3, 1, 2},
		elem{4, 2, 1},
		elem{1, 1, 2},
		elem{5, 1, 2},
	}

	sortElem(e)
	if len(e) != 4 {
		t.Errorf("exected len=4, got %d", len(e))
	}

	if e[0].val != 1 {
		t.Errorf("expected e{1} first, got %#v", e[0])
	}
	if e[1].val != 3 {
		t.Errorf("expected e{3} second, got %#v}", e[1])
	}
}

func rankTest(t *testing.T, e elem, v float64, min, max int) {
	if e.rmin != min || e.rmax != max {
		t.Errorf("expected elem{%2.1f,%d,%d} to be [%2.1f,%d,%d]", e.val, e.rmin, e.rmax, v, min, max)
	}
}

func TestMerge(t *testing.T) {
	return
	l1 := e(0, 1, 2, 3)
	r1 := e()
	ret := merge(l1, r1)
	if len(ret) != 4 {
		t.Errorf("expected len=4; received %#v", ret)
	}
	rankTest(t, ret[0], 0, 1, 1)
	rankTest(t, ret[1], 1, 2, 2)
	rankTest(t, ret[2], 2, 3, 3)
	rankTest(t, ret[3], 3, 4, 4)

	l2 := e()
	r2 := e(0, 1, 2, 3)
	ret = merge(l2, r2)
	if len(ret) != 4 {
		t.Errorf("expected len=4; received %#v", ret)
	}
	rankTest(t, ret[0], 0, 1, 1)
	rankTest(t, ret[1], 1, 2, 2)
	rankTest(t, ret[2], 2, 3, 3)
	rankTest(t, ret[3], 3, 4, 4)

	l3 := e(0, 1, 5, 7)
	r3 := e(1, 1, 2, 9, 11, 15)
	ret = merge(l3, r3)
	if len(ret) != 10 {
		t.Errorf("expected len=8; received %#v", ret)
	}
	rankTest(t, ret[0], 0, 1, 1)
	rankTest(t, ret[1], 1, 2, 2)
	rankTest(t, ret[2], 1, 3, 3)
	rankTest(t, ret[3], 1, 4, 4)
	rankTest(t, ret[4], 2, 5, 5)
	rankTest(t, ret[5], 5, 6, 6)
	rankTest(t, ret[6], 7, 7, 7)
	rankTest(t, ret[7], 9, 8, 8)
	rankTest(t, ret[8], 11, 9, 9)
	rankTest(t, ret[9], 15, 10, 10)

	l4 := []elem{elem{0, 1, 5}, elem{2, 6, 10}}
	r4 := []elem{elem{1, 1, 5}, elem{3, 6, 8}}
	ret = merge(l4, r4)
	if len(ret) != 4 {
		t.Errorf("expected len=3; received %#v", ret)
	}
	rankTest(t, ret[0], 0, 1, 5)
	rankTest(t, ret[1], 1, 6, 10)
	rankTest(t, ret[2], 2, 11, 15)
	rankTest(t, ret[3], 3, 16, 18)
}

func TestCompress(t *testing.T) {
	s := e(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20)
	c := compress(s, 20, 0.001)

	if len(c) != 11 {
		t.Error("expected len=%d, was %d", 11, len(c))
	}
	rankTest(t, c[0], 1, 1, 1)
	rankTest(t, c[1], 2, 2, 2)
	rankTest(t, c[2], 4, 4, 4)
	rankTest(t, c[3], 6, 6, 6)
	rankTest(t, c[4], 8, 8, 8)
	rankTest(t, c[5], 10, 10, 10)
	rankTest(t, c[6], 12, 12, 12)
	rankTest(t, c[7], 14, 14, 14)
	rankTest(t, c[8], 16, 16, 16)
	rankTest(t, c[9], 18, 18, 18)
	rankTest(t, c[10], 20, 20, 20)
}

// test a summary where we end up storing every element
func TestPreciseSummaryStructure(t *testing.T) {
	testSz := 100
	testErr := 0.001

	// test increasing numbers
	s := NewSummaryStructure(testSz, testErr)
	for i := 0; i < testSz; i++ {
		s.Update(float64(20 * i))
	}

	// get each increasing number
	h := s.Histogram()
	errSpan := float64(testSz) * testErr * 20
	for i := 0; i < testSz; i++ {
		quant := h.Quantile(i + 1)
		expect := float64(20 * i)

		if math.Abs(quant-expect) > errSpan {
			t.Errorf("expected %f; received %2.1f, exceeding allowable error of %2.3f [%d]", expect, quant, testErr, int(errSpan))
		}
	}
}

func TestSummaryStructure(t *testing.T) {
	// test nicely round numbers
	testIncreasingSummary(t, 1000, 0.01)
	testIncreasingSummary(t, 750, 0.01)
	testIncreasingSummary(t, 500, 0.01)
	testIncreasingSummary(t, 250, 0.01)
	testIncreasingSummary(t, 100, 0.01)
	testIncreasingSummary(t, 50, 0.01)
	testIncreasingSummary(t, 10, 0.01)
	testIncreasingSummary(t, 5, 0.01)

	// test weird numbers to validate the various floor()
	// and ceil()s that are happening
	// testIncreasingSummary(t, 123, 0.01)
}

func testIncreasingSummary(t *testing.T, testSz int, testErr float64) {
	// test increasing numbers
	s := NewSummaryStructure(testSz, testErr)
	for i := 0; i < testSz; i++ {
		s.Update(float64(i))
	}

	// get each increasing number
	h := s.Histogram()
	errSpan := float64(testSz) * testErr
	for i := 1; i <= testSz; i++ {
		quant := h.Quantile(i)
		expect := float64(i - 1)

		if math.Abs(quant-expect) > errSpan {
			t.Errorf("expected %f; received %2.1f, exceeding allowable error of %2.3f [%d]", expect, quant, testErr, int(errSpan))
		}
	}
}

func BenchmarkSummaryStructure(b *testing.B) {
	testSz := b.N
	testErr := 0.001
	s := NewSummaryStructure(testSz, testErr)

	if testSz < 1000 {
		return
	}

	b.ResetTimer()
	for i := 0; i < testSz; i++ {
		s.Update(float64(i))
	}
}
