package main

import (
	_ "circuit/load"
	fs "circuit/use/anchorfs"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
)

var (
	templates = template.Must(template.ParseGlob("templates/*"))
)

func main() {
	log.Printf("obelisk/")
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("err: %s", err.Error())
	}
}

func error(rw http.ResponseWriter, msg string) {
	log.Printf("err: %s", msg)
	fmt.Fprintf(rw, "Error: %s", msg)
}

func handler(rw http.ResponseWriter, req *http.Request) {
	root := req.URL.Path

	dir, err := fs.OpenDir(root)
	if err != nil {
		error(rw, err.Error())
		return
	}

	dirs, err := dir.Dirs()
	if err != nil {
		error(rw, err.Error())
		return
	}

	rev, workers, err := dir.Files()
	if err != nil {
		error(rw, err.Error())
		return
	}

	fmt.Fprintf(rw, "Obelisk %s -- %s (rev.%d)\n", time.Now().Format(time.Kitchen), root, rev)
	for _, dir := range dirs {
		fmt.Fprintf(rw, "  <a href='%s/%s'>%s/%s/</a>\n", root, dir, root, dir)
	}
	for _, id := range workers {
		fmt.Fprintf(rw, "  %s/%s\n", root, id)
	}
}
