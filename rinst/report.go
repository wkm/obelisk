package rinst

import (
	"fmt"
	"io"
)

// create a textual report of the values
func TextReport(w io.Writer, coll *Collection) {
	b := make(MeasurementBuffer)
	go func() {
		coll.Snapshot(b)
		close(b)
	}()

	select {
	case m := <-b:
		fmt.Fprintf(w, "%25s: %s\n", m.Name, m.Value)
	}
}
