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
		{"\"12\"", "12", ""},
		{"\"hello world\"", "hello world", ""},
		{"\"hello\nworld\"", "hello\nworld", ""},
		{"\"test\\\"escaping\\\"things\" rem", "test\"escaping\"things", " rem"},
	}

	for _, tc := range testcases {
		t.Logf("testcase %q", tc)
		token, rem, err := nextToken(tc.in)

		if err != nil {
			t.Errorf("unexpected error %q", err.Error())
		}

		if token != tc.token {
			t.Errorf("expected token %q, got %q", tc.token, token)
		}

		if rem != tc.rem {
			t.Errorf("expected rem %q, got %q", tc.rem, rem)
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
		{"-123", -123, ""},
		{" 123", 123, ""},
		{" 123  ", 123, "  "},
		{" 123 token", 123, " token"},
	}

	for _, tc := range testcases {
		t.Logf("test case %q", tc)
		r, v, err := parseInt(tc.in)

		if err != nil {
			t.Errorf("unexpected error: %q", err.Error())
		}

		if !v.IsValid() || v.Int() != int64(tc.val) {
			t.Errorf("expected %d, got %q", tc.val, v.Int())
		}

		if r != tc.rem {
			t.Errorf("expected remaining string %q; got %q", tc.rem, r)
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
		{"+12.3", 12.3, ""},
		{"12.3", 12.3, ""},
		{"  12", 12, ""},
		{"12.3 foo", 12.3, " foo"},
	}

	for _, tc := range testcases {
		t.Logf("test case %q", tc)
		r, v, err := parseFloat(tc.in)

		if err != nil {
			t.Errorf("unexpected error: %q", err.Error())
		}

		if !v.IsValid() || math.Abs(v.Float()-tc.val) > 0.001 {
			t.Errorf("expected %d, got %q", tc.val, v.Float())
		}

		if r != tc.rem {
			t.Errorf("expected %q remaining; got %q", tc.rem, r)
		}
	}
}

func TestParseString(t *testing.T) {
	testcases := []struct {
		in, val, rem string
	}{
		{"foo bar", "foo", " bar"},
		{"\"dodge this\" fool", "dodge this", " fool"},
		{"`hi\nthere` fewl", "hi\nthere", " fewl"},
	}

	for _, tc := range testcases {
		t.Logf("test case %q", tc)
		r, v, err := parseString(tc.in)

		if err != nil {
			t.Errorf("unexpected error: %q", err.Error())
		}

		if !v.IsValid() || v.String() != tc.val {
			t.Errorf("expected %q, got %q", tc.val, v.String())
		}

		if r != tc.rem {
			t.Errorf("expected remainder %q; got %q", tc.rem, r)
		}
	}
}
