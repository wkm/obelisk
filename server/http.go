package server

import (
	"bufio"
	"fmt"
	"net/http"

	"github.com/bmizerany/pat"
)

// StartHttpResponder begins listening for HTTP requests on the given address.
// FIXME this should be a HandlerFunc and wired up in the main() func.
func (s *App) StartHttpResponder(httpAddress string) {
	mux := pat.New()
	http.Handle("/", mux)

	mux.Get("/tag/:tag", http.HandlerFunc(s.tagHandler))
	mux.Get("/get/:key", http.HandlerFunc(s.keyHandler))
	mux.Get("/timeseries/:key", http.HandlerFunc(s.timeseriesHandler))

	log.Printf("Starting to listen on %s", httpAddress)
	http.ListenAndServe(httpAddress, nil)
}

func (s *App) tagHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Not implemented\n")
}

func (s *App) keyHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get(":key")
	bb, err := s.kvdb.Get(key)
	bw := bufio.NewWriter(w)

	log.Printf("err=%v, val=%v", err, bb)

	if err != nil {
		fmt.Fprintf(bw, "Error: %s\n", err.Error())
		bw.Flush()
		return
	}

	if len(bb) < 1 {
		fmt.Fprintf(bw, "Unknown Key\n")
		bw.Flush()
		return
	}

	bw.Write(bb)
	bw.Flush()
}

func (s *App) timeseriesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Not implemented\n")
}
