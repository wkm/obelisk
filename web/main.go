package main

import (
	_ "circuit/load/cmd"
	"fmt"
	"github.com/wkm/obelisk/lib/rlog"
	"github.com/wkm/obelisk/server/util"
	"net/http"
	"runtime"
)

var log = rlog.LogConfig.Logger("web")

func main() {
	log.Printf("github.com/wkm/obelisk/")
	log.Printf("setting maxprocs to %d", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())

	xServer, err = util.DiscoverObeliskServer()
	hosts, err := ChildrenTags("host")
	if err != nil {
		log.Printf("could not retrieve list of hosts %s", err.Error())
	}

	log.Printf("hosts: %v", hosts)

	http.HandleFunc("/zk/", zkHandler)
	http.HandleFunc("/host/", hostHandler)
	http.HandleFunc("/worker/", workerHandler)
	http.HandleFunc("/service/", serviceHandler)

	http.HandleFunc("/api/time", timeHandler)

	http.HandleFunc("/", indexHandler)

	log.Printf("starting HTTP")
	err = http.ListenAndServe(":8080", LogAccess(http.DefaultServeMux))
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

func indexHandler(rw http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	if path == "/" {
		path = "/index.html"
	}
	// try to find a template with the given name
	renderTemplate(req, rw, path, nil)
}
