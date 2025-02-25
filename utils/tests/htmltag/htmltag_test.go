package htmltag

import (
	"testing"
)

func TestContainsAll(t *testing.T) {
	tests := []struct {
		items  []string
		target []string
		want   bool
	}{
		{[]string{"apple", "banana"}, []string{"banana", "orange", "apple", "grape"}, true},
		{[]string{"apple", "banana", "cherry"}, []string{"banana", "orange", "apple", "grape"}, false},
		{[]string{}, []string{"banana", "orange", "apple", "grape"}, true},
		{[]string{"banana"}, []string{}, false},
		{[]string{"apple", "banana"}, []string{"apple", "banana"}, true},
	}

	for _, tt := range tests {
		got := containsAll(tt.items, tt.target)
		if got != tt.want {
			t.Errorf("containsAll(%v, %v) = %v; want %v", tt.items, tt.target, got, tt.want)
		}
	}
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
				{"div > Content", ""},
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
				{"div > attrs[id]", "main"},
				{"div > attrs[class]", "container"},
				{"div > attrs[data-custom]", "value"},
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
				{"div > attrs[id]", "parent"},
				{"#parent > p.child > Content", "Text"},
				{"#parent > span > attrs[data-test]", "true"},
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
				{"div > attrs[id]", "level1"},
				{"#level2 > attrs[id]", "level2"},
				{"#level3 > attrs[id]", "level3"},
				{"#level3 > p.deep > Content", "Content"},
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
				{"section > attrs[class]", "outer"},
				{"section > article > attrs[data-type]", "news"},
				{"section > article > div.content > span.highlight > Content", "Text"},
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
				{"div > attrs[id]", "root"},
				{"#root > section > attrs[class]", "main"},
				{"#root > section > header > h1 > Content", "Title"},
				{"#root > section > footer > p.footer-text > Content", "Footer"},
			},
			wantErr: false,
		},
		{
			name: "Multiclass acces",
			html: `<div id="root"><section class="main"><header><h1>Title</h1></header><footer><p class="footer-text last">Footer</p></footer></section></div>`,
			pathAssertions: []struct {
				path          string
				expectedValue string
			}{
				{"#root > section.main > footer > p.footer-text.last > Content", "Footer"},
			},
			wantErr: false,
		},

		{
			name: "Simple text in header",
			html: `<h1>Title</h1>`,
			pathAssertions: []struct {
				path          string
				expectedValue string
			}{
				{"h1 > Content", "Title"},
			},
			wantErr: false,
		},
		{
			name: "Nested HTML with mixed content",
			html: `<section><p>Text <strong>bold</strong> text.</p></section>`,
			pathAssertions: []struct {
				path          string
				expectedValue string
			}{
				{"section > p > Content", "Text bold text."},
				{"section > p > strong > Content", "bold"},
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
				AssertPath(t, got, assertion.path, assertion.expectedValue)
			}
		})
	}
}
