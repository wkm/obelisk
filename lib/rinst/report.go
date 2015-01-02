package rinst

import (
	"fmt"
	"io"
)

// TextReport creates a human readable report of the values in plaintext.
func TextReport(w io.Writer, inst Instrument) {
	b := make(MeasurementBuffer, 0)
	inst.Measure("", &b)

	for _, m := range b {
		fmt.Fprintf(w, "%s: %d %f\n", m.Name, m.IntValue, m.FloatValue)
	}
}
