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
	children, stat, _ := zk.Children(root)
	nodes := make(map[string]*NodeInfo)
	for _, child := range children {
		childNode := path.Join(root, child)
		log.Printf("childnode: %s", childNode)

		stat, err := zk.Exists(childNode)
		ni := new(NodeInfo)
		nodes[child] = ni

		log.Printf(" -- info: %s", ni)

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
		log.Printf(" -- datalen: %d", datalen)
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

	var parent string
	if root != "/" {
		parent = ".."
	}

	zk := zkResponse{parent, root, stat, nodes}
	log.Printf("ZK: %s", zk)
	renderTemplate(rw, "/zk.html", &zk)
}
