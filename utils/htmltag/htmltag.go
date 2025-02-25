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
}

func HtmlToTag(htmlStr string) (*Tag, error) {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return nil, err
	}

	firstNode := doc.FirstChild.FirstChild.NextSibling.FirstChild // Navigate to <div>

	if firstNode == nil {
		return nil, fmt.Errorf("empty node")
	}

	return nodeToTag(firstNode), nil
}

func nodeToTag(n *html.Node) *Tag {
	if n.Type == html.ElementNode {
		tag := &Tag{
			TagName: n.Data,
			Attrs:   make(map[string]string),
		}
		for _, attr := range n.Attr {
			tag.Attrs[attr.Key] = attr.Val
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			childTag := nodeToTag(c)
			if childTag != nil {
				tag.Children = append(tag.Children, *childTag)
			}
		}
		return tag
	}
	return nil
}
