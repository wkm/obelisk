package rinst

import (
	"fmt"
	"io"
)

// create a textual report of the values
func TextReport(w io.Writer, inst Instrument) {
	b := make(MeasurementBuffer)
	go func() {
		inst.Measure("", b)
		close(b)
	}()

	for {
		select {
		case m, ok := <-b:
			if !ok {
				return
			}

			fmt.Fprintf(w, "%s: %s\n", m.Name, m.Value)
		}
	}
}
