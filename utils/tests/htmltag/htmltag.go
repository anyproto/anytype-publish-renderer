package htmltag

import (
	"fmt"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
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
	firstNode := doc.FirstChild.FirstChild.NextSibling

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

func UnequalFail(t assert.TestingT, expected, actual string, doDiff bool, msgAndArgs ...interface{}) bool {
	diff := "<diff disabled>"
	if doDiff {
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(expected, actual, false)
		diff = dmp.DiffPrettyText(diffs)
	}

	return assert.Fail(t, fmt.Sprintf("Not equal: \n"+
		"expected: %s\n"+
		"actual  : %s\n"+
		"diff    : %s", expected, actual, diff), msgAndArgs...)
}

// assertPath checks if the given path exists in the Tag structure and matches the expected value.
// Path looks like this: "section.main > footer > p.footer-text.last > Content"
//
// Note:
// Unlike CSS accessors, it has to contain _all_ nodes, not just some of them.
// I.e., it doesn't traverse all children nodes recursively.
func AssertPath(t assert.TestingT, tag *Tag, path string, expectedValue string) bool {

	if tag == nil {
		return UnequalFail(t, "<tag>", "nil", false, "Expected a Tag, but got nil")
	}

	path = strings.ReplaceAll(path, " >", ">")
	path = strings.ReplaceAll(path, "> ", ">")
	parts := strings.Split(path, ">")
	current := tag
	for i, part := range parts {
		if strings.HasPrefix(part, "#") && i < len(parts)-1 {
			// Handle ID selection
			id := strings.TrimPrefix(parts[i], "#")
			current = findTagById(current, id)
			if current == nil {
				expected := fmt.Sprintf("`%s`", id)
				msg := fmt.Sprintf("Expected to find element with id %s, but it does not exist", id)
				return UnequalFail(t, expected, "`nil`", false, msg)
			}
		} else if i == len(parts)-1 {
			// Last part, check if it's an attribute, content or tag name
			if strings.Contains(part, "attrs[") {
				// Attribute access, e.g., "attrs[id]"
				attrName := strings.TrimPrefix(strings.TrimSuffix(part, "]"), "attrs[")
				if current.Attrs[attrName] != expectedValue {
					expected := fmt.Sprintf("`%s`", expectedValue)
					actual := fmt.Sprintf("`%s`", current.Attrs[attrName])
					msg := fmt.Sprintf("Expected attribute %s to be %s, but got %s", attrName, expectedValue, current.Attrs[attrName])
					return UnequalFail(t, expected, actual, true, msg)
				}
			} else if part == "Content" {
				// Content access
				if current.Content != expectedValue {
					expected := fmt.Sprintf("`%s`", expectedValue)
					actual := fmt.Sprintf("`%s`", current.Content)
					msg := fmt.Sprintf("Expected content to be %s, but got %s", expectedValue, current.Content)
					return UnequalFail(t, expected, actual, true, msg)
				}
			} else {
				// Tag name access
				if current.TagName != expectedValue {
					expected := fmt.Sprintf("`%s`", expectedValue)
					actual := fmt.Sprintf("`%s`", current.TagName)
					msg := fmt.Sprintf("Expected tag name to be %s, but got %s", expectedValue, current.TagName)
					return UnequalFail(t, expected, actual, true, msg)
				}
			}
		} else {
			classAccessors := strings.Split(part, ".")
			classes := make([]string, 0)
			tagName := part

			// Check classes if they are present in selector
			if len(classAccessors) > 1 {
				tagName = classAccessors[0]
				classes = classAccessors[1:]
			}
			// Navigate to the child tag
			found := false
			for _, child := range current.Children {
				childClasses := strings.Fields(child.Attrs["class"])
				if child.TagName == tagName && containsAll(classes, childClasses) {
					current = &child
					found = true
					break
				}
			}
			if !found {
				expected := fmt.Sprintf("`%s`", part)
				actual := "`nil`"
				msg := fmt.Sprintf("Expected to find child tag %s, but it does not exist", part)
				return UnequalFail(t, expected, actual, false, msg)
			}
		}
	}
	return true
}

// containsAll checks that all `items` are present in `target`
func containsAll(items []string, target []string) bool {
	targetSet := make(map[string]struct{})

	// Populate the target set
	for _, t := range target {
		targetSet[t] = struct{}{}
	}

	// Check if all items exist in the target set
	for _, item := range items {
		if _, found := targetSet[item]; !found {
			return false
		}
	}

	return true
}

// findTagById searches for a tag with the specified ID within the current tag and its descendants.
func findTagById(tag *Tag, id string) *Tag {
	if tag == nil {
		return nil
	}
	if tag.Attrs["id"] == id {
		return tag
	}
	for _, child := range tag.Children {
		if result := findTagById(&child, id); result != nil {
			return result
		}
	}
	return nil
}
