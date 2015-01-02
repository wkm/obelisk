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
			// "Pools":    "/pool",
			"Hosts": "/host",
			// "Services": "/service",
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

// HeaderBar prints the header, optionally with an item highlighted
func HeaderBar(names map[string]string, active string) string {
	node := NewHTMLNode("ul").AddClass("left")
	divider := NewHTMLNode("li").AddClass("divider")
	node.AddChild(divider)

	for name, link := range names {
		section := NewHTMLNode("li").AddClass("name")
		link := NewHTMLNode("a").AddAttribute("href", link)
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

// string print a bunch of numbers comma separated
func commaSeparated(items []uint64) string {
	itemStrs := make([]string, len(items))
	lastPositive := 0
	for i, obj := range items {
		if obj > 0 {
			itemStrs[i] = fmt.Sprintf("%d", obj)
			lastPositive = i
		}
	}
	return strings.Join(itemStrs[:lastPositive], ",")
}
