package internal

import (
	"fmt"
	"sync"
)

func downloadWorker(wg *sync.WaitGroup, linksToProcess <-chan string) {
	defer wg.Done()

	for link := range linksToProcess {
		fileName := GetFileNameFromUrl(link)
		fmt.Printf("Downloading %s\n", fileName)
		err := DownloadFileFromUrl(link, fileName)

		if err != nil {
			fmt.Printf("Error downloading file %s %v", link, err)
		}
	}
}

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

	fmt.Printf("Found %d valid image links\n", len(links))

	numberOfWorkers := len(links)

	linksToProcess := make(chan string)

	var wg sync.WaitGroup

	for i := 0; i < numberOfWorkers; i++ {
		wg.Add(1)
		go downloadWorker(&wg, linksToProcess)
	}

	for _, link := range links {
		linksToProcess <- link
	}

	close(linksToProcess)

	wg.Wait()

	return nil
}
