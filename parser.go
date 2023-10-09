package crawler

import (
	"io"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

func dfs(node *html.Node, actor func(*html.Node)) {
	actor(node)
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		dfs(child, actor)
	}
}

func parseLinks(html_reader io.Reader) ([]Link, error) {
	root, err := html.Parse(html_reader)
	if err != nil {
		return nil, err
	}

	links := shortArray[Link]()

	dfs(root, func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					links = append(links, Link{Url: strings.TrimSpace(attr.Val), Text: strings.TrimSpace(extractText(node))})
					break
				}
			}
		}
	})

	return links, nil
}

func extractText(node *html.Node) string {
	text := ""
	dfs(node, func(node *html.Node) {
		if node.Type == html.TextNode {
			text += strings.TrimSpace(node.Data)
		}
	})

	return strings.ReplaceAll(text, "\n", " ")
}

func parsePage(html_reader io.Reader) (string, string, error) {
	root, err := html.Parse(html_reader)
	if err != nil {
		return "", "", err
	}

	var title, content string
	dfs(root, func(node *html.Node) {
		if node.Type == html.TextNode && node.Parent != nil && isHtmlPageContent(node.Parent) {
			content += strings.TrimSpace(node.Data) + " "
		} else if node.Type == html.TextNode && node.Parent != nil && isHtmlTitle(node.Parent) {
			title = strings.TrimSpace(node.Data)
		}
	})

	content = strings.TrimSpace(content)

	return title, content, nil
}

func isHtmlTitle(node *html.Node) bool {
	return node.Type == html.ElementNode && node.Data == "title" && node.Parent != nil && node.Parent.Data == "head"
}

var htmlContentTags = []string{
	"p", "h1", "h2", "h3", "h4", "h5", "h6",
}

func isHtmlPageContent(node *html.Node) bool {
	return node.Type == html.ElementNode && slices.Contains(htmlContentTags, node.Data)
}
