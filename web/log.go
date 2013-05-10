package main

import (
	"net/http"
	"obelisk/lib/rlog"
)

var accessLog = rlog.LogConfig.Logger("access")

func LogAccess(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessLog.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
