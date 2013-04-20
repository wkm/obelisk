package main

import (
	"net/http"
	// "text/template"
	"circuit/kit/zookeeper"
	"circuit/kit/zookeeper/zutil"
	"log"
	"strings"
)

type zkResponse struct {
	Parent string
	Node   string
	Stat   *zookeeper.Stat
	Dirs   []string
}

var (
	zk, err = zutil.DialUntilReady("127.0.0.1:2181")
)

func zkHandler(rw http.ResponseWriter, req *http.Request) {
	root := req.URL.Path[4:]
	root = strings.TrimSuffix(root, "/")

	if !strings.HasPrefix(root, "/") {
		root = "/" + root
	}

	log.Printf("zk handler: %s", root)
	children, stat, _ := zk.Children(root)

	var parent string
	if root != "/" {
		parent = ".."
	}

	zk := zkResponse{parent, root, stat, children}
	log.Printf("ZK: %s", zk)
	renderTemplate(rw, "/zk.html", &zk)
}
