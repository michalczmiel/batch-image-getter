package provider

import (
	"io"
	"regexp"
	"strings"

	"github.com/michalczmiel/batch-image-getter/internal"
	"golang.org/x/net/html"
)

type HtmlProvider struct {
	url         string
	httpClient  internal.HttpClient
	parameters  *internal.Parameters
	regexSearch bool
}

func NewHtmlProvider(url string, httpClient internal.HttpClient, parameters *internal.Parameters, regexSearch bool) Provider {
	return &HtmlProvider{
		url:         url,
		httpClient:  httpClient,
		parameters:  parameters,
		regexSearch: regexSearch,
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

func (p *HtmlProvider) findImageLinksUsingRegex(content string) []string {
	var links []string

	urlPatterns := []*regexp.Regexp{
		regexp.MustCompile(`https?://[^\s<>"']+`),
		regexp.MustCompile(`www\.[^\s<>"']+`),
		regexp.MustCompile(`href=["']([^"']+)["']`),
		regexp.MustCompile(`(?:href)=["'](/[^"']+)["']`),
	}

	for _, pattern := range urlPatterns {
		matches := pattern.FindAllString(content, -1)

		for _, match := range matches {
			if strings.HasSuffix(match, "/") || strings.HasSuffix(match, "\\") {
				match = match[:len(match)-1]
			}

			if !internal.IsUrlValid(match) {
				continue
			}

			isValidImage := false
			for _, allowedExtension := range p.parameters.ImageTypes {
				if strings.HasSuffix(match, allowedExtension) {
					isValidImage = true
					break
				}
			}

			if !isValidImage {
				continue
			}

			links = append(links, match)
		}
	}

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

	var rawLinks []string

	if p.regexSearch {
		content, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		rawLinks = p.findImageLinksUsingRegex(string(content))
	} else {
		doc, err := html.Parse(response.Body)
		if err != nil {
			return nil, err
		}
		rawLinks = GetImageLinksFromHtmlDoc(doc)
	}

	return internal.ProcessLinks(p.url, rawLinks), nil
}
