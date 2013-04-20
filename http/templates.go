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
	err := getTemplates().ExecuteTemplate(rw, name, object)
	if err != nil {
		respondError(rw, err.Error())
	}
}

func dirLastModified(dir string) time.Time {
	var t int64 = 0
	filepath.Walk(dir, func(path string, io os.FileInfo, err error) error {
		this := io.ModTime().Unix()
		if this > t {
			t = this
		}
		return nil
	})

	return time.Unix(t, 0)
}

func getTemplates() *template.Template {
	lastModified := dirLastModified(templateDir)

	if lastModified.Unix() > lastUpdated.Unix() {
		lastUpdated = lastModified
		log.Printf("building templates from %s", lastUpdated.Format(time.Kitchen))
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
