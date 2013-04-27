package timestore

import (
	"fmt"
	"github.com/petar/GoLLRB/llrb"
	"io"
	"time"
)

// write out all values into a writer, this blocks the entire time
func (s *Store) Dump(w io.Writer) error {
	statDump.Incr()

	s.Lock()
	defer s.Unlock()

	var cnt int
	var err error
	hasErrored := false
	iterfn := func(i llrb.Item) bool {
		// don't iterate if we've errored
		if hasErrored {
			return false
		}

		p := i.(Point)
		cnt, err = fmt.Fprintf(w, "%d=%f,", p.Time, p.Value)
		statIter.Add(uint(cnt))
		if err != nil {
			statError.Incr()
			hasErrored = true
		}

		return true
	}

	// iterate over all keys
	for key, series := range s.values {
		count, err := fmt.Fprintf(w, "%d:", key)
		statDumpBytes.Add(uint(count))
		if err != nil {
			statError.Incr()
			return err
		}

		series.AscendGreaterOrEqual(Point{0, 0}, iterfn)
		if hasErrored {
			return err
		}

		count, err = fmt.Fprintf(w, "\n")
		statDumpBytes.Add(uint(count))
		if err != nil {
			statError.Incr()
			return err
		}
	}

	return nil
}

// execute a flush periodically
func (s *Store) startPeriodicFlush() {
	switch {
	case <-s.shutdown:
		return

	case time <-s.flushTick; time != nil:
		statFlush.Incr()
		time.Sleep(s.config.FlushPeriod)
		s.Flush()
	}
}

func (s *Store) Flush() {

}
