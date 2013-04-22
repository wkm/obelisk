package main

import (
	"net/http"
	// "text/template"
	"circuit/kit/zookeeper"
	"circuit/kit/zookeeper/zutil"
	"circuit/load/config"
	"circuit/sys/zanchorfs"
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
	Anchor *zanchorfs.ZFile
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
	zk, err = zutil.DialUntilReady(config.Config.Zookeeper.Zookeepers())
	dataMax = 2048
)

func zkHandler(rw http.ResponseWriter, req *http.Request) {
	root := req.URL.Path
	root = strings.TrimPrefix(root, "/zk/")
	root = strings.TrimSuffix(root, "/")

	zkr := new(zkResponse)

	if !strings.HasPrefix(root, "/") {
		root = "/" + root
	}

	zkr.Node = root
	zkr.Stat, err = zk.Exists(root)
	if err != nil {
		respondError(rw, err.Error())
		return
	}

	if zkr.Stat == nil {
		respondError(rw, "Unknown node")
		return
	}

	if zkr.Stat.DataLength() > 0 {
		zkr.Data, _, err = zk.Get(root)
		if err != nil {
			respondError(rw, err.Error())
		}
	}

	afile, err := getAnchorFile(zkr.Data)
	if err != nil {
		log.Printf("not an anchor file")
	} else {
		log.Printf("yes, an anchor file: %s", afile)
		zkr.Anchor = afile
	}

	if root != "/" {
		zkr.Parent = ".."
	}

	// get data on children nodes
	children, _, err := zk.Children(root)
	if err != nil {
		respondError(rw, err.Error())
	}

	nodes := make(map[string]*NodeInfo)
	for _, child := range children {
		childNode := path.Join(root, child)

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

	zkr.Nodes = nodes
	renderTemplate(rw, "/zk.html", zkr)
}
