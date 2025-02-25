package htmltag

import (
	"bytes"

	"golang.org/x/net/html"
)

// Tag structure to be used in tests to replace html string comparison
type Tag struct {
	TagName  string
	Attrs    map[string]string
	Children []Tag
}

func HtmlToTag(htmlStr string) (*Tag, error) {
	doc, err := html.Parse(bytes.NewReader([]byte(htmlStr)))
	if err != nil {
		return nil, err
	}
	return nodeToTag(doc), nil
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
