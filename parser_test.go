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

	expectEqualInTest(t, extractText(node), "abc")
}

func TestExtractTextReplaceNewlineToSpace(t *testing.T) {
	node, err := html.Parse(strings.NewReader("<p>a\nb</p>"))
	if err != nil {
		t.Fatal(err)
	}

	expectEqualInTest(t, extractText(node), "a b")
}

func TestExtractTextReturnsTextOfAnchor(t *testing.T) {
	node, err := html.Parse(strings.NewReader(`
		<a href="example.com">text</a>
	`))
	if err != nil {
		t.Fatal(err)
	}

	expectEqualInTest(t, extractText(node), "text")
}

func TestParseLinksReturnsLinksWithText(t *testing.T) {
	links, err := parseLinks(strings.NewReader(`
		<a href="example.com">text</a>
	`))
	if err != nil {
		t.Fatal(err)
	}

	expectEqualInTest(t, links, []Link{{Url: "example.com", Text: "text"}})
}

func TestParseLinksTrimsWhitespaces(t *testing.T) {
	links, err := parseLinks(strings.NewReader(`
		<a href=" example.com "> text </a>
	`))
	if err != nil {
		t.Fatal(err)
	}

	expectEqualInTest(t, links, []Link{{Url: "example.com", Text: "text"}})
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

	expectEqualInTest(t, links, []Link{
		{Url: "a.com", Text: "a"},
		{Url: "b.com", Text: "b"},
	})
}

func TestParsePageReturnsTitleAndContent(t *testing.T) {
	title, content, err := parsePage(strings.NewReader(`
		<html>
			<head>
				<title>Title</title>
			</head>
			<body>
				<h1>Header</h1>
				<p>Content</p>
			</body>
		</html>
	`))
	if err != nil {
		t.Fatal(err)
	}

	expectEqualInTest(t, title, "Title")
	expectEqualInTest(t, content, "Header Content")
}

func TestParsePageIgnoresNonContentTags(t *testing.T) {
	_, content, err := parsePage(strings.NewReader(`
		<html>
			<body>
				<span>span</span>
				<button>Click Me!</button>
			</body>
		</html>
	`))
	if err != nil {
		t.Fatal(err)
	}

	expectEqualInTest(t, content, "")

}
