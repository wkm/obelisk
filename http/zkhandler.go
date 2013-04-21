package main

import (
	"net/http"
	// "text/template"
	"circuit/kit/zookeeper"
	"circuit/kit/zookeeper/zutil"
	"log"
	"path"
	"strings"
)

type zkResponse struct {
	Parent string
	Node   string
	Stat   *zookeeper.Stat
	Nodes  map[string]*NodeInfo
}

type NodeInfo struct {
	Name  string
	IsDir bool
	Stat  *zookeeper.Stat
	Data  string
	Error string
}

var (
	zk, err = zutil.DialUntilReady("127.0.0.1:2181")
	dataMax = 2048
)

func zkHandler(rw http.ResponseWriter, req *http.Request) {
	root := req.URL.Path[4:]
	root = strings.TrimSuffix(root, "/")

	if !strings.HasPrefix(root, "/") {
		root = "/" + root
	}

	log.Printf("root: " + root)
	children, stat, _ := zk.Children(root)
	nodes := make(map[string]*NodeInfo)
	for _, child := range children {
		childNode := path.Join(root, child)
		log.Printf("childnode: " + childNode)

		stat, err := zk.Exists(childNode)
		ni := new(NodeInfo)
		nodes[child] = ni

		ni.Name = child
		ni.Stat = stat
		if err != nil {
			log.Printf("error: " + err.Error())
			ni.Error = err.Error()
			continue
		}

		if stat.NumChildren() > 0 {
			ni.IsDir = true
		}

		if stat.DataLength() < dataMax {
			data, _, err := zk.Get(childNode)
			if err != nil {
				log.Printf("error: " + err.Error())
				ni.Error = err.Error()
				continue
			}
			ni.Data = data
		}
	}

	var parent string
	if root != "/" {
		parent = ".."
	}

	zk := zkResponse{parent, root, stat, nodes}
	log.Printf("ZK: %s", zk)
	renderTemplate(rw, "/zk.html", &zk)
}
