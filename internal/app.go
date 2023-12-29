package internal

import (
	"fmt"
	"sync"
)

func DownloadImagesFromWebsite(url string, imageTypesToDownload []string) error {
	doc, err := GetHtmlDocFromUrl(url)
	if err != nil {
		return err
	}

	rawLinks := GetImageLinksFromHtmlDoc(doc)
	if len(rawLinks) == 0 {
		return fmt.Errorf("no links found")
	}

	links := ProcessLinks(url, rawLinks, imageTypesToDownload)

	var wg sync.WaitGroup
	for _, link := range links {
		wg.Add(1)

		go func(l string) {
			defer wg.Done()

			fileName := GetFileNameFromUrl(l)
			err = DownloadFileFromUrl(l, fileName)
			if err != nil {
				fmt.Printf("Error downloading file %s %v", l, err)
			}
		}(link)
	}

	wg.Wait()

	return nil
}
