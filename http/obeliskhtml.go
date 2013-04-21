package main

import (
	"bytes"
)

var (
	header = func() map[string]string {
		return map[string]string{
			"Hosts":   "/host.html",
			"Workers": "/workers.html",
			"ZK":      "/zk",
		}
	}

	htmlHelpers = map[string]interface{}{
		"headerBar": HeaderBar,
		"header":    header,
	}
)

func HeaderBar(names map[string]string, active string) string {
	node := NewHtmlNode("ul").AddClass("left")
	divider := NewHtmlNode("li").AddClass("divider")
	node.AddChild(divider)

	for name, link := range names {
		section := NewHtmlNode("li").AddClass("name")
		link := NewHtmlNode("a").AddAttribute("href", link)
		link.AddChild(NewTextNode(name))
		section.AddChild(link)

		if name == active {
			section.AddClass("active")
		}

		node.AddChild(section)
		node.AddChild(divider)
	}

	var doc bytes.Buffer
	node.Write(&doc)
	return doc.String()
}
