package main

import (
	"bytes"
	"fmt"
	// "reflect"
	"strings"
)

var (
	header = func() map[string]string {
		return map[string]string{
			"Hosts":   "/host",
			"Workers": "/worker",
			"ZK":      "/zk",
		}
	}

	htmlHelpers = map[string]interface{}{
		"headerBar":      HeaderBar,
		"header":         header,
		"CommaSeparated": commaSeparated,
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

func commaSeparated(items []uint64) string {
	itemStrs := make([]string, len(items))
	for i, obj := range items {
		if obj < 1 {
			itemStrs[i] = fmt.Sprintf("%d", obj)
		}
	}
	return strings.Join(itemStrs, ",")
}
