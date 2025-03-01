package provider

import (
	"reflect"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func docWithLinks(t *testing.T) *html.Node {
	r := strings.NewReader(`
		<html>
			<head>
				<title>Test</title>
			</head>
			<body>
				<a href="https://example.com">Example</a>
				<div>
					<img src="https://example.com/image.png" />
					<img src="https://example.com/logo.svg" />
				</div>
				<span data-src="//example.com/image3.png"/>
				<img src="/image2.png" data-src="https://example.com/image2.png" />
			</body>
		</html>
	`)

	doc, err := html.Parse(r)
	if err != nil {
		t.Fatal(err)
	}

	return doc
}

func TestGetImageLinksFromHtmlDoc(t *testing.T) {
	expectedLinks := []string{
		"https://example.com/image.png",
		"https://example.com/logo.svg",
		"//example.com/image3.png",
		"/image2.png",
		"https://example.com/image2.png",
	}

	links := GetImageLinksFromHtmlDoc(docWithLinks(t))

	if !reflect.DeepEqual(links, expectedLinks) {
		t.Fatalf("Expected %v links, got %v", expectedLinks, links)
	}
}

// TODO: add test for regex search
