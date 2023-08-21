package crawler

import (
	"errors"
	"io"
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

func parseTitle(html_reader io.Reader) (string, error) {
	root, err := html.Parse(html_reader)
	if err != nil {
		return "", err
	}

	title := ""
	dfs(root, func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "title" && node.Parent != nil && node.Parent.Data == "head" && node.FirstChild != nil {
			title = node.FirstChild.Data
		}
	})

	if title == "" {
		return "", errors.New("Did not find title")
	}

	return title, nil
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
