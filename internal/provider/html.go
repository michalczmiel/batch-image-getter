package provider

import (
	"github.com/michalczmiel/batch-image-getter/internal"
	"golang.org/x/net/html"
)

type HtmlProvider struct {
	url        string
	httpClient internal.HttpClient
	parameters *internal.Parameters
}

func NewHtmlProvider(url string, httpClient internal.HttpClient, parameters *internal.Parameters) Provider {
	return &HtmlProvider{
		url:        url,
		httpClient: httpClient,
		parameters: parameters,
	}
}

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

func (p *HtmlProvider) Links() ([]string, error) {
	response, err := p.httpClient.Request(p.url, map[string]string{
		"Referer": p.parameters.Referer,
	})
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	doc, err := html.Parse(response.Body)
	if err != nil {
		return nil, err
	}

	rawLinks := GetImageLinksFromHtmlDoc(doc)

	return internal.ProcessLinks(p.url, rawLinks), nil
}
