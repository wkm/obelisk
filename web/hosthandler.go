package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
)

type HostInfo struct {
	Name    string
	Metrics []*MetricInfo
}

type MetricInfo struct {
	Name string
	Path string
}

func hostHandler(rw http.ResponseWriter, req *http.Request) {
	root := req.URL.Path
	root = strings.TrimPrefix(root, "/host/")

	// list all hosts
	if root == "" {
		children, err := ChildrenTags("host")
		if err != nil {
			respondError(rw, err.Error())
			return
		}

		renderTemplate(req, rw, "/allhosts.html", children)
		return
	}

	log.Printf("host data for %s", root)

	tags, err := ChildrenTags("host", root)
	if err != nil {
		respondError(rw, err.Error())
		return
	}
	sort.Strings(tags)

	metrics := make([]*MetricInfo, len(tags))
	for i, tag := range tags {
		m := new(MetricInfo)
		m.Name = tag
		m.Path = filepath.Join("host", root, tag)
		metrics[i] = m
	}
	renderTemplate(req, rw, "/host.html", &HostInfo{root, metrics})
}
