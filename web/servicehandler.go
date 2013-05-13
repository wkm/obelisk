package main

import (
	"net/http"
)

func serviceHandler(rw http.ResponseWriter, req *http.Request) {
	renderTemplate(req, rw, "service.html", nil)
}
