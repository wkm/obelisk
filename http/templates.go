package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

var (
	templateDir     = "templates"
	lastUpdated     = time.Unix(0, 0)
	cachedTemplates *template.Template
)

func renderTemplate(rw http.ResponseWriter, name string, object interface{}) {
	err := getTemplates().ExecuteTemplate(rw, name, nil)
	if err != nil {
		error(rw, err.Error())
	}
}

func getTemplates() *template.Template {
	f, err := os.Open(templateDir)
	if err != nil {
		log.Printf("could not find template directory %s", err.Error())
		return nil
	}

	fi, err := f.Stat()
	if err != nil {
		log.Printf("could not stat template directory %s", err.Error())
	}

	if fi.ModTime().Unix() > lastUpdated.Unix() {
		log.Printf("building templates")
		lastUpdated = fi.ModTime()
		cachedTemplates = buildTemplates("templates")
	} else {
		// use cached
	}

	return cachedTemplates
}

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
