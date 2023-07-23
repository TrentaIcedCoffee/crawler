package crawler

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestExtractTextReturnsTextWithOrder(t *testing.T) {
	node, err := html.Parse(strings.NewReader(`
		<p>a</p>
		<p>b</p>
		<p>c</p>
	`))
	if err != nil {
		t.Fatal(err)
	}

	expectEqual(t, extractText(node), "abc")
}

func TestExtractTextReplaceNewlineToSpace(t *testing.T) {
	node, err := html.Parse(strings.NewReader("<p>a\nb</p>"))
	if err != nil {
		t.Fatal(err)
	}

	expectEqual(t, extractText(node), "a b")
}

func TestExtractTextReturnsTextOfAnchor(t *testing.T) {
	node, err := html.Parse(strings.NewReader(`
		<a href="example.com">text</a>
	`))
	if err != nil {
		t.Fatal(err)
	}

	expectEqual(t, extractText(node), "text")
}

func TestParseLinksReturnsLinksWithText(t *testing.T) {
	links, err := parseLinks(strings.NewReader(`
		<a href="example.com">text</a>
	`))
	if err != nil {
		t.Fatal(err)
	}

	expectEqual(t, links, []Link{{"example.com", "text"}})
}

func TestParseLinksReturnsMultipleLinkWithText(t *testing.T) {
	links, err := parseLinks(strings.NewReader(`
		<div>
			<a href="a.com">a</a>
			<div>
				<a href="b.com">b</a>
			</div>
		</div>
	`))
	if err != nil {
		t.Fatal(err)
	}

	expectEqual(t, links, []Link{
		{"a.com", "a"},
		{"b.com", "b"},
	})
}
