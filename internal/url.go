package internal

import (
	"net/url"
	"path"
	"strings"
)

func GetFileNameFromUrl(imageUrl string) (string, error) {
	u, err := url.Parse(imageUrl)
	if err != nil {
		return "", err
	}

	return path.Base(u.Path), nil
}

func IsUrlValid(rawUrl string) bool {
	_, err := url.ParseRequestURI(rawUrl)

	return err == nil
}

func resolveUrl(url, link string) string {
	if strings.HasPrefix(link, "//") {
		// is protocol relative link
		return "https:" + link
	} else if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
		return link
	} else {
		// is relative link
		return url + link
	}
}

func RemoveDuplicates(original []string) []string {
	// set is not available in Go, so we use map instead to remove duplicates
	var withoutDuplicatesSet = map[string]struct{}{}

	for _, item := range original {
		withoutDuplicatesSet[item] = struct{}{}
	}

	var output []string
	for item := range withoutDuplicatesSet {
		output = append(output, item)
	}

	return output
}

func ProcessLinks(url string, rawLinks []string) []string {
	var processedLinks []string

	for _, link := range rawLinks {
		var fullUrl = resolveUrl(url, link)

		processedLinks = append(processedLinks, fullUrl)
	}

	return processedLinks
}
