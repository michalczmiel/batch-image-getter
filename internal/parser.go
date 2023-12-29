package internal

import (
	"strings"

	"golang.org/x/net/html"
)

func isValidImageLink(link string, imageTypes []string) bool {
	for _, suffix := range imageTypes {
		if strings.HasSuffix(link, suffix) {
			return true
		}
	}

	return false
}

func isNodeWithImageLink(n *html.Node) bool {
	return n.Type == html.ElementNode && (n.Data == "img" || n.Data == "span")
}

func GetImageLinksFromHtmlDoc(doc *html.Node, imageTypes []string) []string {
	var links []string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if isNodeWithImageLink(n) {
			for _, attr := range n.Attr {
				if isValidImageLink(attr.Val, imageTypes) {
					links = append(links, attr.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return links
}
