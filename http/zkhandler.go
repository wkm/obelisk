package main

import (
	"net/http"
	// "text/template"
	"circuit/kit/zookeeper"
	"circuit/kit/zookeeper/zutil"
	"log"
	"path"
	"strings"
	"time"
)

type zkResponse struct {
	Parent string
	Node   string
	Data   string
	Stat   *zookeeper.Stat
	Nodes  map[string]*NodeInfo
}

type NodeInfo struct {
	Name     string
	IsDir    bool
	Stat     *zookeeper.Stat
	Data     string
	Error    string
	MTimeStr string
}

var (
	zk, err = zutil.DialUntilReady("zk1.datacenter.net:2181")
	dataMax = 2048
)

func zkHandler(rw http.ResponseWriter, req *http.Request) {
	root := req.URL.Path[4:]
	root = strings.TrimSuffix(root, "/")

	if !strings.HasPrefix(root, "/") {
		root = "/" + root
	}

	log.Printf("root: " + root)
	stat, err := zk.Exists(root)
	if err != nil {
		respondError(rw, err.Error())
	}

	var data string
	if stat.DataLength() > 0 {
		data, _, err = zk.Get(root)
		if err != nil {
			respondError(rw, err.Error())
		}
	}

	var parent string
	if root != "/" {
		parent = ".."
	}

	// get data on children nodes

	log.Printf("children ...")
	children, stat, err := zk.Children(root)
	if err != nil {
		respondError(rw, err.Error())
	}

	nodes := make(map[string]*NodeInfo)
	for _, child := range children {
		childNode := path.Join(root, child)
		log.Printf("childnode: %s", childNode)

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

		ni.MTimeStr = stat.MTime().Format(time.RFC822)

		if stat.NumChildren() > 0 {
			ni.IsDir = true
		}

		datalen := stat.DataLength()
		if datalen > 0 && datalen < dataMax {
			data, _, err := zk.Get(childNode)
			if err != nil {
				log.Printf("error: " + err.Error())
				ni.Error = err.Error()
				continue
			}
			ni.Data = data
		}
	}

	zk := zkResponse{parent, root, data, stat, nodes}
	renderTemplate(rw, "/zk.html", &zk)
}
