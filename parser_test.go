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

	ExpectEqualInTest(t, extractText(node), "abc")
}

func TestExtractTextReplaceNewlineToSpace(t *testing.T) {
	node, err := html.Parse(strings.NewReader("<p>a\nb</p>"))
	if err != nil {
		t.Fatal(err)
	}

	ExpectEqualInTest(t, extractText(node), "a b")
}

func TestExtractTextReturnsTextOfAnchor(t *testing.T) {
	node, err := html.Parse(strings.NewReader(`
		<a href="example.com">text</a>
	`))
	if err != nil {
		t.Fatal(err)
	}

	ExpectEqualInTest(t, extractText(node), "text")
}

func TestParseLinksReturnsLinksWithText(t *testing.T) {
	links, err := parseLinks(strings.NewReader(`
		<a href="example.com">text</a>
	`))
	if err != nil {
		t.Fatal(err)
	}

	ExpectEqualInTest(t, links, []Link{{Url: "example.com", Text: "text"}})
}

func TestParseLinksTrimsWhitespaces(t *testing.T) {
	links, err := parseLinks(strings.NewReader(`
		<a href=" example.com "> text </a>
	`))
	if err != nil {
		t.Fatal(err)
	}

	ExpectEqualInTest(t, links, []Link{{Url: "example.com", Text: "text"}})
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

	ExpectEqualInTest(t, links, []Link{
		{Url: "a.com", Text: "a"},
		{Url: "b.com", Text: "b"},
	})
}

func TestParseTitleReturnsTitle(t *testing.T) {
	title, err := parseTitle(strings.NewReader(`
		<html>
			<head>
				<title>Title</title>
			</head>
		</html>
	`))
	if err != nil {
		t.Fatal(err)
	}

	ExpectEqualInTest(t, title, "Title")
}

func TestParseTitleReturnsErrorIfAbsent(t *testing.T) {
	title, err := parseTitle(strings.NewReader(`
		<html>
			<head>
			</head>
		</html>
	`))
	if err == nil {
		t.Fatalf("Expect to error when title is absent, but found %s as title", title)
	}
}

func TestParseTitleReturnsErrorIfEmptyTitle(t *testing.T) {
	title, err := parseTitle(strings.NewReader(`
		<html>
			<head>
				<title></title>
			</head>
		</html>
	`))
	if err == nil {
		t.Fatalf("Expect to error when title is absent, but found %s as title", title)
	}
}
