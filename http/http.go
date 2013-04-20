package main

import (
	_ "circuit/load"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"
)

// var (
// 	templates = buildTemplates("templates")
// )

func buildTemplates(directory string) *template.Template {
	files, err := filepath.Glob(directory + "/*")
	if err != nil {
		log.Printf("err: %s", err.Error())
	}

	t := template.New("")

	log.Printf("Files: %s", files)
	for _, file := range files {
		// skip partials
		// if strings.HasPrefix(filepath.Base(file), "_") {
		// 	continue
		// }

		name := strings.Replace(file, directory, "", 1)
		t = t.New(name)
		art, err := ioutil.ReadFile(file)
		if err != nil {
			log.Printf("err: %s", err.Error())
			continue
		}

		t, err = t.Parse(string(art))
		if err != nil {
			log.Printf("err: %s", err.Error())
		}

		t.ParseFiles(file)

		log.Printf("created template `%s`", name)
	}

	return t
}

func main() {
	log.Printf("obelisk/")
	http.HandleFunc("/zk/", zkHandler)
	http.HandleFunc("/", indexHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("err: %s", err)
	}
}

func getTemplates() *template.Template {
	return buildTemplates("templates")
}

func error(rw http.ResponseWriter, msg string) {
	rw.Header().Set("Status-Code", "502")
	log.Printf("err: %s", msg)
	err := getTemplates().ExecuteTemplate(rw, "error.html", msg)
	if err != nil {
		fmt.Fprintf(rw, "Could not render error template: %s", err.Error())
	}
}

func renderTemplate(rw http.ResponseWriter, name string, object interface{}) {
	err := getTemplates().ExecuteTemplate(rw, name, nil)
	if err != nil {
		error(rw, err.Error())
	}
}

func indexHandler(rw http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	if path == "/" {
		path = "/index.html"
	}
	// try to find a template with the given name
	err := getTemplates().ExecuteTemplate(rw, path, nil)
	if err != nil {
		log.Printf("err: %s", err.Error())
		error(rw, err.Error())
	}
}
