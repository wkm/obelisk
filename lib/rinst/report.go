package rinst

import (
	"fmt"
	"io"
)

// create a textual report of the values
func TextReport(w io.Writer, inst Instrument) {
	b := make(MeasurementBuffer, 0)
	inst.Measure("", &b)

	for _, m := range b {
		fmt.Fprintf(w, "%s: %d %f\n", m.Name, m.IntValue, m.FloatValue)
	}
}
