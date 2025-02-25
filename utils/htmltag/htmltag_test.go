package htmltag

import (
	"strings"
	"testing"
)

// assertPath checks if the given path exists in the Tag structure and matches the expected value.
func assertPath(t *testing.T, tag *Tag, path string, expectedValue string) {
	parts := strings.Split(path, ".")
	current := tag
	i := 0

	for i < len(parts) {
		part := parts[i]
		if part == "#" && i < len(parts)-1 {
			// Handle ID selection
			id := parts[i+1]
			current = findTagById(current, id)
			if current == nil {
				t.Errorf("Expected to find element with id %s, but it does not exist", id)
				return
			}
			// Skip the next part since it's the ID
			i += 2
		} else if i == len(parts)-1 {
			// Last part, check if it's an attribute or tag name
			if strings.Contains(part, "[") {
				// Attribute access, e.g., "attrs[id]"
				attrName := strings.TrimPrefix(strings.TrimSuffix(part, "]"), "[")
				if current.Attrs[attrName] != expectedValue {
					t.Errorf("Expected attribute %s to be %s, but got %s", attrName, expectedValue, current.Attrs[attrName])
				}
			} else {
				// Tag name access
				if current.TagName != expectedValue {
					t.Errorf("Expected tag name to be %s, but got %s", expectedValue, current.TagName)
				}
			}
			i++
		} else {
			// Navigate to the child tag
			found := false
			for _, child := range current.Children {
				if child.TagName == part {
					current = &child
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected to find child tag %s, but it does not exist", part)
				return
			}
			i++
		}
	}
}

// findTagById searches for a tag with the specified ID within the current tag and its descendants.
func findTagById(tag *Tag, id string) *Tag {
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

func TestHtmlToTag(t *testing.T) {
	tests := []struct {
		name           string
		html           string
		pathAssertions []struct {
			path          string
			expectedValue string
		}
		wantErr bool
	}{
		{
			name: "Simple HTML with one element",
			html: "<div></div>",
			pathAssertions: []struct {
				path          string
				expectedValue string
			}{
				{"TagName", "div"},
			},
			wantErr: false,
		},
		{
			name: "HTML with multiple attributes",
			html: `<div id="main" class="container" data-custom="value"></div>`,
			pathAssertions: []struct {
				path          string
				expectedValue string
			}{
				{"attrs[id]", "main"},
				{"attrs[class]", "container"},
				{"attrs[data-custom]", "value"},
			},
			wantErr: false,
		},
		{
			name: "Nested HTML elements with attributes",
			html: `<div id="parent"><p class="child">Text</p><span data-test="true"></span></div>`,
			pathAssertions: []struct {
				path          string
				expectedValue string
			}{
				{"attrs[id]", "parent"},
				{"p.attrs[class]", "child"},
				{"span.attrs[data-test]", "true"},
			},
			wantErr: false,
		},
		{
			name: "Deeply nested HTML elements",
			html: `<div id="level1"><div id="level2"><div id="level3"><p class="deep">Content</p></div></div></div>`,
			pathAssertions: []struct {
				path          string
				expectedValue string
			}{
				{"attrs[id]", "level1"},
				{"#level2.attrs[id]", "level2"},
				{"#level3.attrs[id]", "level3"},
				{"#level3.p.attrs[class]", "deep"},
			},
			wantErr: false,
		},
		{
			name: "Nested HTML with mixed attributes",
			html: `<section class="outer"><article data-type="news"><div class="content"><span class="highlight">Text</span></div></article></section>`,
			pathAssertions: []struct {
				path          string
				expectedValue string
			}{
				{"attrs[class]", "outer"},
				{"article.attrs[data-type]", "news"},
				{"article.div.attrs[class]", "content"},
				{"article.div.span.attrs[class]", "highlight"},
			},
			wantErr: false,
		},
		{
			name: "Complex nested structure",
			html: `<div id="root"><section class="main"><header><h1>Title</h1></header><footer><p class="footer-text">Footer</p></footer></section></div>`,
			pathAssertions: []struct {
				path          string
				expectedValue string
			}{
				{"attrs[id]", "root"},
				{"section.attrs[class]", "main"},
				{"section.header.h1.TagName", "h1"},
				{"#root.section.footer.p.attrs[class]", "footer-text"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HtmlToTag(tt.html)
			if (err != nil) != tt.wantErr {
				t.Errorf("HtmlToTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, assertion := range tt.pathAssertions {
				assertPath(t, got, assertion.path, assertion.expectedValue)
			}
		})
	}
}
