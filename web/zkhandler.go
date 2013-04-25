package main

import (
	"circuit/kit/zookeeper"
	"circuit/kit/zookeeper/zutil"
	"circuit/load/config"
	"circuit/sys/zanchorfs"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
)

type zkResponse struct {
	Parent  string
	Node    string
	Data    string
	GobData string
	Stat    *zookeeper.Stat
	Anchor  *zanchorfs.ZFile
	Nodes   map[string]*NodeInfo
}

type NodeInfo struct {
	Name     string
	IsDir    bool
	Stat     *zookeeper.Stat
	Data     string
	GobData  string
	Error    string
	MTimeStr string
}

var (
	zk, err = zutil.DialUntilReady(config.Config.Zookeeper.Zookeepers())
	dataMax = 2048
)

func zkHandler(rw http.ResponseWriter, req *http.Request) {
	root := req.URL.Path

	// have to redirect to a node ending in '/' for relative links
	if !strings.HasSuffix(root, "/") {
		redirectTo(rw, req, req.URL.Path+"/")
		return
	}

	root = strings.TrimPrefix(root, "/zk")
	root = strings.TrimSuffix(root, "/")
	if root == "" {
		root = "/"
	}

	switch req.URL.RawQuery {
	case "delete":
		log.Printf("Deleting node %s", root)
		err := zk.Delete(root, -1)
		if err != nil {
			respondError(rw, err.Error())
			return
		}

		setFlash(rw, fmt.Sprintf("deleted node %s", root))
		redirectTo(rw, req, path.Join(req.URL.Path, ".."))

	case "deletechildren":
		log.Printf("Delete children of node %s", root)
		children, _, err := zk.Children(root)
		if err != nil {
			respondError(rw, err.Error())
			return
		}

		for _, child := range children {
			childnode := path.Join(root, child)
			log.Printf(" -- deleting %s", childnode)
			err := zk.Delete(childnode, -1)
			if err != nil {
				respondError(rw, err.Error())
				return
			}
		}

		setFlash(rw, fmt.Sprintf("deleted %d children of node %s", len(children), root))
		redirectTo(rw, req, req.URL.Path)

	default:
		zkr := new(zkResponse)

		if root != "/" {
			zkr.Parent = ".."
		}

		log.Printf("root node: %s", root)
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
			// not an anchor file
		} else {
			log.Printf("yes, an anchor file: %s", afile)
			zkr.Anchor = afile
		}

		str, err := getAsGob(zkr.Data)
		if err == nil {
			log.Printf("yes, gob: %s", str)
			zkr.GobData = str
		} else {
			log.Printf("not gob because %s", err.Error())
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

				str, err := getAsGob(ni.Data)
				if err == nil {
					ni.GobData = str
				}
			}
		}

		zkr.Nodes = nodes
		renderTemplate(req, rw, "/zk.html", zkr)
	}
}
