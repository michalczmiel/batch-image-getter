package internal

import (
	"golang.org/x/net/html"
)

type Provider interface {
	Links() ([]string, error)
}

type HtmlProvider struct {
	url        string
	httpClient HttpClient
	parameters *Parameters
}

func NewHtmlProvider(url string, httpClient HttpClient, parameters *Parameters) Provider {
	return &HtmlProvider{
		url:        url,
		httpClient: httpClient,
		parameters: parameters,
	}
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

	return ProcessLinks(p.url, rawLinks), nil
}

type FileProvider struct {
	path       string
	fileSystem FileSystem
}

func NewFileProvider(path string, fileSystem FileSystem) Provider {
	return &FileProvider{
		path:       path,
		fileSystem: fileSystem,
	}
}

func (p *FileProvider) Links() ([]string, error) {
	lines, err := p.fileSystem.ReadLines(p.path)
	if err != nil {
		return nil, err
	}

	var links []string
	for _, line := range lines {
		if IsUrlValid(line) {
			links = append(links, line)
		}
	}

	return links, nil
}
