package main

import (
	"net/http"

	"github.com/wkm/obelisk/lib/rlog"
)

var accessLog = rlog.LogConfig.Logger("access")

// LogAccess writes HTTP requests and response metadata into the log.
func LogAccess(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessLog.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		webRequest.Incr()
		// var start = time.Now().UnixNano()
		handler.ServeHTTP(w, r)
		// webRespTime.Add(time.Now().UnixNano() - start)
	})
}
