package main

import (
	"bytes"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// Global variables enabling template refreshing.
var (
	templateDir     = "templates"
	lastUpdated     = time.Unix(0, 0)
	cachedTemplates *template.Template
	cachedFiles     = make(map[string]*[]byte)
)

// TemplateResponse is the top-level object fed templates containing information
// about the request.
type TemplateResponse struct {
	Req    *http.Request
	Flash  string
	Object interface{}
}

// renderTemplate into an HTTP response, with the right Content-Type.
func renderTemplate(req *http.Request, rw http.ResponseWriter, name string, object interface{}) {
	flash := getFlash(req)
	deleteFlash(rw)

	name = strings.TrimPrefix(name, "/")
	ext := filepath.Ext(name)
	rw.Header().Set("Content-Type", mime.TypeByExtension(ext))

	// check if we have cached the file
	data, ok := cachedFiles[name]
	if ok {
		rw.Write(*data)
		return
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

// dirLastModified gives the newest last modified inode within all levels of a directory.
func dirLastModified(dir string) time.Time {
	var t int64
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

// getTemplates gives the list of templates, rebuilding them if the directory has changed.
func getTemplates() *template.Template {
	lastModified := dirLastModified(templateDir)

	if lastModified.Unix() > lastUpdated.Unix() {
		lastUpdated = lastModified
		log.Printf("building templates and buffers from %s", lastUpdated.Format(time.Kitchen))
		cachedFiles = make(map[string]*[]byte)
		cachedTemplates = buildTemplates("templates")
	} else {
		// use cached
	}

	return cachedTemplates
}

// buildTemplates recursively builds all templates in a directory.
func buildTemplates(directory string) *template.Template {
	log.Printf("building templates")
	t := template.New("")
	var templateWalkFn func(path string, io os.FileInfo, err error) error
	templateWalkFn = func(path string, io os.FileInfo, err error) error {
		// recurse into directories
		if io.IsDir() {
			return nil
		}

		// non-html is buffered
		if filepath.Ext(path) != ".html" {
			buildBuffer(directory, path)
		} else {
			buildTemplate(directory, path, t)
		}

		// process template file
		return nil
	}

	err := filepath.Walk(directory, templateWalkFn)

	if err != nil {
		log.Printf("err: %s", err.Error())
	}

	return t
}

func buildBuffer(directory, path string) {
	name := strings.TrimPrefix(strings.Replace(path, directory, "", 1), "/")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("err: %s", err.Error())
		return
	}

	cachedFiles[name] = &data
	log.Printf("buffer `%s`", name)
}

func buildTemplate(directory, path string, t *template.Template) {
	name := strings.TrimPrefix(strings.Replace(path, directory, "", 1), "/")
	t = t.New(name)
	t.Funcs(htmlHelpers)
	art, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("err: %s", err.Error())
		return
	}

	t, err = t.Parse(string(art))
	if err != nil {
		log.Printf("err: %s", err.Error())
	}

	log.Printf("template `%s`", name)
}
