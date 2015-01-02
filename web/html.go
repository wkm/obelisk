package main

import (
	"fmt"
	"io"
	"strings"
	"text/template"
)

// HTMLNode is a naive representation of an HTML dom node
type HTMLNode struct {
	Tag        string
	Class      []string
	Attributes map[string]string
	Children   []HTMLWriter
}

// RawHTMLNode represents pre-formatted HTML content.
type RawHTMLNode string

// HTMLWriter represents a type which can write formatted HTML.
type HTMLWriter interface {
	Write(r io.Writer) (int, error)
}

// TextNode is a simple textual node.
type TextNode string

// NewHTMLNode creates a new HTML node with the given tag.
func NewHTMLNode(tag string) *HTMLNode {
	n := new(HTMLNode)
	n.Tag = tag
	n.Class = make([]string, 0)
	n.Attributes = make(map[string]string, 0)
	n.Children = make([]HTMLWriter, 0)
	return n
}

// AddAttribute adds an attribute to the HTML tag.
func (n *HTMLNode) AddAttribute(name, value string) *HTMLNode {
	n.Attributes[name] = value
	return n
}

// AddClass adds classes to the HTML tag.
func (n *HTMLNode) AddClass(name ...string) *HTMLNode {
	n.Class = append(n.Class, name...)
	return n
}

// AddChild adds subnodes into the HTML node.
func (n *HTMLNode) AddChild(node ...HTMLWriter) *HTMLNode {
	n.Children = append(n.Children, node...)
	return n
}

// Write prints the node as formatted HTML into the given writer.
func (n *HTMLNode) Write(r io.Writer) (int, error) {
	// create the basic tag and class attribute
	count, err := fmt.Fprintf(r, "<%s class='%s'\n", n.Tag, strings.Join(n.Class, " "))
	if err != nil {
		return count, err
	}

	// write out other attributes
	for key, value := range n.Attributes {
		incr, err := fmt.Fprintf(r, " %s='%s'", key, value)
		count += incr
		if err != nil {
			return count, err
		}
	}

	// close out the open tag
	incr, err := fmt.Fprintf(r, ">")
	count += incr
	if err != nil {
		return count, err
	}

	// write out child nodes
	for _, child := range n.Children {
		incr, err := child.Write(r)
		count += incr
		if err != nil {
			return count, err
		}
	}

	// write out the close tag
	incr, err = fmt.Fprintf(r, "</%s>\n", n.Tag)
	count += incr
	if err != nil {
		return count, err
	}
	return count, nil
}

// Write prints the preformatted HTML without further interpretation.
func (n *RawHTMLNode) Write(w io.Writer) (int, error) {
	return fmt.Fprint(w, string(*n))
}

// NewTextNode creates a new HTML node which contains raw text.
func NewTextNode(s string) *TextNode {
	node := TextNode(s)
	return &node
}

// Write prints escaped text into the HTML node.
func (t *TextNode) Write(w io.Writer) (int, error) {
	template.HTMLEscape(w, []byte(*t))
	return 0, nil
}
