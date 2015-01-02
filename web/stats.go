package main

import (
	"github.com/wkm/obelisk/lib/rinst"
)

var (
	// WebStats contains metrics on the web service
	WebStats   = rinst.NewCollection()
	webRequest = WebStats.Counter("reqs", "req", "http requests received")
	// webRespSize = WebStats.Distribution("resp.sz", "byte", "size of http web responses")
	// webRespTime = WebStats.Distribution("resp.time", "ns","response time to web requests")
)
