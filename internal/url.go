package internal

import (
	"net/url"
	"strings"
)

func GetFileNameFromUrl(url string) string {
	urlParts := strings.Split(url, "/")

	return urlParts[len(urlParts)-1]
}

func IsUrlValid(rawUrl string) bool {
	_, err := url.ParseRequestURI(rawUrl)

	return err == nil
}

func isValidImageLink(link string, imageTypes []string) bool {
	for _, suffix := range imageTypes {
		// // sometimes image links have uppercase extensions
		if strings.HasSuffix(strings.ToLower(link), suffix) {
			return true
		}
	}

	return false
}

func ProcessLinks(url string, rawLinks, fileTypes []string) []string {
	// set is not available in Go, so we use map instead to remove duplicates
	var validLinks = map[string]struct{}{}

	for _, link := range rawLinks {
		if !isValidImageLink(link, fileTypes) {
			continue
		}

		var fullUrl string
		if strings.HasPrefix(link, "//") {
			// is protocol relative link
			fullUrl = "https:" + link
		} else if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
			fullUrl = link
		} else {
			// is relative link
			fullUrl = url + link
		}

		validLinks[fullUrl] = struct{}{}
	}

	var links []string
	for link := range validLinks {
		links = append(links, link)
	}
	return links
}
