package internal

import (
	"golang.org/x/net/html"
)

func isNodeWithImageLink(n *html.Node) bool {
	return n.Type == html.ElementNode && (n.Data == "img" || n.Data == "span")
}

func GetImageLinksFromHtmlDoc(doc *html.Node) []string {
	var links []string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if isNodeWithImageLink(n) {
			for _, attr := range n.Attr {
				if attr.Key == "src" || attr.Key == "data-src" {
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
