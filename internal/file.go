package internal

import (
	"strings"
)

func GetFileNameFromUrl(url string) string {
	urlParts := strings.Split(url, "/")

	return urlParts[len(urlParts)-1]
}
