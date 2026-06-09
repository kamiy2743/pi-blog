package article

import (
	"bytes"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func ConvertMarkdownToHTML(markdown string) (string, error) {
	var rendered bytes.Buffer
	md := goldmark.New(
		goldmark.WithExtensions(extension.Table),
	)
	if err := md.Convert([]byte(markdown), &rendered); err != nil {
		return "", err
	}

	sanitized := bluemonday.UGCPolicy().Sanitize(rendered.String())
	return openLinksInBlankTab(sanitized)
}

func openLinksInBlankTab(htmlFragment string) (string, error) {
	contextNode := &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Body,
		Data:     "body",
	}
	nodes, err := html.ParseFragment(strings.NewReader(htmlFragment), contextNode)
	if err != nil {
		return "", err
	}

	for _, node := range nodes {
		addBlankTabAttrsToLinks(node)
	}

	var rendered bytes.Buffer
	for _, node := range nodes {
		if err := html.Render(&rendered, node); err != nil {
			return "", err
		}
	}
	return rendered.String(), nil
}

func addBlankTabAttrsToLinks(node *html.Node) {
	if node.Type == html.ElementNode && node.Data == "a" {
		setAttr(node, "target", "_blank")
		setAttr(node, "rel", "noopener noreferrer")
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		addBlankTabAttrsToLinks(child)
	}
}

func setAttr(node *html.Node, key string, value string) {
	for i := range node.Attr {
		if node.Attr[i].Key == key {
			node.Attr[i].Val = value
			return
		}
	}

	node.Attr = append(node.Attr, html.Attribute{
		Key: key,
		Val: value,
	})
}
