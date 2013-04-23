package main

import (
	"bytes"
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

// a template response is the top-level object fed templates containing information
// about the request
type TemplateResponse struct {
	Req    *http.Request
	Flash  string
	Object interface{}
}

// render a template into an HTTP response, with the right Content-Type
func renderTemplate(req *http.Request, rw http.ResponseWriter, name string, object interface{}) {
	flash := getFlash(req)
	deleteFlash(rw)

	name = strings.TrimPrefix(name, "/")
	ext := filepath.Ext(name)

	// FIXME use mime.GetTypeByExtension
	switch ext {
	case ".html":
		rw.Header().Set("Content-Type", "text/html")

	case ".css":
		rw.Header().Set("Content-Type", "text/css")
	}

	// we have to buffer response writes so we have an opportunity to modify the headers
	// (eg, set flashes, etc.)
	var buff bytes.Buffer
	err := getTemplates().ExecuteTemplate(&buff, name, &TemplateResponse{req, flash, object})
	if err != nil {
		respondError(rw, err.Error())
	} else {
		buff.WriteTo(rw)
	}
}

// get the newest last modified inode within all levels of a directory
func dirLastModified(dir string) time.Time {
	var t int64 = 0
	filepath.Walk(dir, func(path string, io os.FileInfo, err error) error {
		this := io.ModTime().Unix()
		if this > t {
			t = this
		}

		// keep going
		return nil
	})

	return time.Unix(t, 0)
}

// get the templates, rebuilding them if the directory has changed
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

// recursively build all templates in a directory
func buildTemplates(directory string) *template.Template {
	log.Printf("building templates")
	t := template.New("")
	var templateWalkFn func(path string, io os.FileInfo, err error) error
	templateWalkFn = func(path string, io os.FileInfo, err error) error {
		// process template file
		name := strings.TrimPrefix(strings.Replace(path, directory, "", 1), "/")
		t = t.New(name)
		t.Funcs(htmlHelpers)
		art, err := ioutil.ReadFile(path)
		if err != nil {
			log.Printf("err: %s", err.Error())
			return nil
		}

		t, err = t.Parse(string(art))
		if err != nil {
			log.Printf("err: %s", err.Error())
		}

		log.Printf("created template `%s`", name)
		return nil
	}

	err := filepath.Walk(directory, templateWalkFn)

	if err != nil {
		log.Printf("err: %s", err.Error())
	}

	return t
}
