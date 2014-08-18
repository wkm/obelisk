package main

import (
	"net/http"
	"path/filepath"
	"sort"
	"strings"

	"github.com/wkm/obelisk/lib/rinst"
)

type HostInfo struct {
	Name    string
	Metrics map[string]*MetricGroup
	Workers []string
}

type MetricGroup struct {
	Name string
	Info map[string]*MetricInfo
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

	workers, err := ChildrenTags("host", root, "workers")
	if err != nil {
		respondError(rw, err.Error())
		return
	}

	groups, err := workerResponseData(workers[0])
	if err != nil {
		respondError(rw, err.Error())
		return
	}

	renderTemplate(req, rw, "/host.html", &HostInfo{root, groups, workers})
}

func workerResponseData(worker string) (map[string]*MetricGroup, error) {
	tags, err := ChildrenTags("worker", worker)
	if err != nil {
		return nil, err
	}
	sort.Strings(tags)

	groups := make(map[string]*MetricGroup)
	for _, tag := range tags {
		paths := strings.SplitN(tag, ".", 2)
		groupName := paths[0]
		if group, ok := groups[groupName]; !ok {
			group = &MetricGroup{groupName, make(map[string]*MetricInfo)}
			groups[groupName] = group
		}

		name := filepath.Join("worker", worker, tag)
		info, err := GetMetricInfo(name)
		if err != nil {
			log.Printf("err: %s: %s", name, err.Error())
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

		case rinst.TypeIntValue, rinst.TypeBoolValue, rinst.TypeFloatValue, rinst.TypeDateValue:
			m.TypeName = "value"
			m.IsRate = false
		}
		m.Path = filepath.Join("worker", worker, tag)
		groups[groupName].Info[m.Path] = m
	}

	return groups, nil
}
