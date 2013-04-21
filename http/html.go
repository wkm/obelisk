package main

import (
	"fmt"
	"io"
	"strings"
	"text/template"
)

// a naive representation of an HTML dom node
type HtmlNode struct {
	Tag        string
	Class      []string
	Attributes map[string]string
	Children   []HtmlWriter
}

type RawHtmlNode string

type HtmlWriter interface {
	Write(r io.Writer) (int, error)
}

type TextNode string

func NewHtmlNode(tag string) *HtmlNode {
	n := new(HtmlNode)
	n.Tag = tag
	n.Class = make([]string, 0)
	n.Attributes = make(map[string]string, 0)
	n.Children = make([]HtmlWriter, 0)
	return n
}

// set an attribute
func (n *HtmlNode) AddAttribute(name, value string) *HtmlNode {
	n.Attributes[name] = value
	return n
}

// add CSS classes
func (n *HtmlNode) AddClass(name ...string) *HtmlNode {
	n.Class = append(n.Class, name...)
	return n
}

// add a child node
func (n *HtmlNode) AddChild(node ...HtmlWriter) *HtmlNode {
	n.Children = append(n.Children, node...)
	return n
}

// write HTML nodes
func (n *HtmlNode) Write(r io.Writer) (int, error) {
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

func (n *RawHtmlNode) Write(w io.Writer) (int, error) {
	return fmt.Fprint(w, string(*n))
}

func NewTextNode(s string) *TextNode {
	node := TextNode(s)
	return &node
}

func (t *TextNode) Write(w io.Writer) (int, error) {
	template.HTMLEscape(w, []byte(*t))
	return 0, nil
}
