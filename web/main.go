package main

import (
	_ "circuit/load"
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Printf("obelisk/")
	http.HandleFunc("/zk/", zkHandler)
	http.HandleFunc("/host/", hostHandler)
	http.HandleFunc("/worker/", workerHandler)
	http.HandleFunc("/", indexHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("err: %s", err)
	}
}

func respondError(rw http.ResponseWriter, msg string) {
	// FIXME doubt this works
	rw.Header().Set("Status-Code", "502")
	log.Printf("err: %s", msg)

	var resp TemplateResponse
	resp.Object = msg

	err := getTemplates().ExecuteTemplate(rw, "error.html", &resp)
	if err != nil {
		fmt.Fprintf(rw, "Could not render error template: %s", err.Error())
	}
}

func redirectTo(rw http.ResponseWriter, req *http.Request, path string) {
	log.Printf("Redirect -> %s", path)
	http.Redirect(rw, req, path, http.StatusFound)
}

func hostHandler(rw http.ResponseWriter, req *http.Request) {
	renderTemplate(req, rw, "/host.html", nil)
}

func indexHandler(rw http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	if path == "/" {
		path = "/index.html"
	}
	// try to find a template with the given name
	renderTemplate(req, rw, path, nil)
}
