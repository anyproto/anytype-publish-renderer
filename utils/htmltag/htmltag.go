package htmltag

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type Tag struct {
	TagName  string
	Attrs    map[string]string
	Children []Tag
	Content  string
}

func HtmlToTag(htmlStr string) (*Tag, error) {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return nil, err
	}

	// html.Parse creates normalized html, with html>head>body structure
	// From the other hand, html.ParseFragment also requires this normalized structure
	// to be passed, which also looks not appealing.
	// Therefore, if err != nil it will always have this first nodes, which we skip to
	// navigate to the content.
	firstNode := doc.FirstChild.FirstChild.NextSibling.FirstChild

	if firstNode == nil {
		return nil, fmt.Errorf("empty node")
	}

	return nodeToTag(firstNode), nil
}

func nodeToTag(n *html.Node) *Tag {
	switch n.Type {
	case html.TextNode:
		return &Tag{Content: strings.TrimSpace(n.Data)}
	case html.ElementNode:
		tag := &Tag{
			TagName: n.Data,
			Attrs:   make(map[string]string),
		}
		for _, attr := range n.Attr {
			tag.Attrs[attr.Key] = attr.Val
		}
		var contentBuilder strings.Builder
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			childTag := nodeToTag(c)
			if c.Type == html.TextNode {
				contentBuilder.WriteString(c.Data)
			} else if childTag != nil {
				tag.Children = append(tag.Children, *childTag)
				if childTag.Content != "" {
					contentBuilder.WriteString(childTag.Content)
				}
			}
		}
		tag.Content = strings.TrimSpace(contentBuilder.String())
		return tag
	}

	return nil
}
