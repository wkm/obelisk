package main

import (
	"github.com/wkm/obelisk/lib/rlog"
	"net/http"
	// "time"
)

var accessLog = rlog.LogConfig.Logger("access")

func LogAccess(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessLog.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		webRequest.Incr()
		// var start = time.Now().UnixNano()
		handler.ServeHTTP(w, r)
		// webRespTime.Add(time.Now().UnixNano() - start)
	})
}
