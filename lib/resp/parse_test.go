package resp

import (
	"math"
	"testing"
)

func TestNextToken(t *testing.T) {
	testcases := []struct {
		in, token, rem string
	}{
		{"12", "12", ""},
		{" 12", "12", ""},
		{" 12 ", "12", " "},
		{" 12 abc", "12", " abc"},
	}

	for _, tc := range testcases {
		t.Logf("testcase %v", tc)
		token, rem, err := nextToken(tc.in)

		if err != nil {
			t.Errorf("unexpected error %v", err.Error())
		}

		if token != tc.token {
			t.Errorf("expected token %v, got %v", tc.token, token)
		}

		if rem != tc.rem {
			t.Errorf("expected rem %v, got %v", tc.rem, rem)
		}
	}
}

func TestParseInt(t *testing.T) {
	testcases := []struct {
		in  string
		val int
		rem string
	}{
		{"123", 123, ""},
		{" 123", 123, ""},
		{" 123  ", 123, "  "},
		{" 123 token", 123, " token"},
	}

	for _, tc := range testcases {
		t.Logf("test case %v", tc)
		r, v, err := parseInt(tc.in)

		if err != nil {
			t.Errorf("unexpected error: %v", err.Error())
		}

		if !v.IsValid() || v.Int() != int64(tc.val) {
			t.Errorf("expected %d, got %v", tc.val, v.Int())
		}

		if r != tc.rem {
			t.Errorf("expected remaining string %v; got %v", tc.rem, r)
		}
	}
}

func TestParseFloat(t *testing.T) {
	testcases := []struct {
		in  string
		val float64
		rem string
	}{
		{"12", 12, ""},
		{"12.3", 12.3, ""},
		{"  12", 12, ""},
	}

	for _, tc := range testcases {
		t.Logf("test case %v", tc)
		r, v, err := parseFloat(tc.in)

		if err != nil {
			t.Errorf("unexpected error: %v", err.Error())
		}

		if !v.IsValid() || math.Abs(v.Float()-tc.val) > 0.001 {
			t.Errorf("expected %d, got %v", tc.val, v.Float())
		}

		if r != tc.rem {
			t.Errorf("expected %v remaining; got %v", tc.rem, r)
		}
	}
}

func TestParseSlice(t *testing.T) {
	testcases := []struct {
		in, val string
	}{
		{"a b cs", "{\"a\", \"b\", \"c\"}"},
		{"a b c\r\nb c d", "{\"a\", \"b\", \"c\"}"},
	}

	for _, tc := range testcases {
		t.Logf("test case %v", tc)
		r, v, err := parseSlice(tc.in)
		t.Logf("parsed into: r=%#v v=%#v err=%#v", r, v, err)
		t.Logf(" len=%d", v.Len())
		for i := 0; i < v.Len(); i++ {
			t.Logf("  %d=%s\t%#v", i, v.Index(i).String(), v.Index(i))
		}
	}
}
