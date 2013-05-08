package main

import (
	"log"
	"net/http"
	"obelisk/lib/rinst"
	"path/filepath"
	"sort"
	"strings"
)

type HostInfo struct {
	Name    string
	Metrics []*MetricInfo
}

type MetricInfo struct {
	Name       string
	Unit, Desc string
	TypeName   string
	IsRate     bool
	Path       string
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
		name := filepath.Join("host", root, tag)
		info, err := GetMetricInfo(name)
		if err != nil {
			log.Printf("err: %s", err.Error())
			continue
		}

		m := new(MetricInfo)
		m.Name = info.Name
		m.Desc = info.Description
		m.Unit = info.Unit
		switch info.Type {
		case rinst.TypeCounter:
			m.TypeName = "counter"
			m.IsRate = true
		case rinst.TypeValue:
			m.TypeName = "value"
			m.IsRate = false
		}
		m.Path = filepath.Join("host", root, tag)
		metrics[i] = m
	}
	renderTemplate(req, rw, "/host.html", &HostInfo{root, metrics})
}
