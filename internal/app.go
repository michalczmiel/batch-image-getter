package internal

import (
	"fmt"
	"sync"
)

func downloadWorker(wg *sync.WaitGroup, linksToProcess <-chan string, failedLinks chan<- error) {
	defer wg.Done()

	for link := range linksToProcess {
		fileName := GetFileNameFromUrl(link)
		fmt.Printf("Downloading %s\n", link)
		err := DownloadFileFromUrl(link, fileName)

		if err != nil {
			failedLinks <- fmt.Errorf("error downloading file %s %v", link, err)
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
	failedLinks := make(chan error, numberOfWorkers)

	var wg sync.WaitGroup

	for i := 0; i < numberOfWorkers; i++ {
		wg.Add(1)
		go downloadWorker(&wg, linksToProcess, failedLinks)
	}

	for _, link := range links {
		linksToProcess <- link
	}
	close(linksToProcess)

	wg.Wait()

	close(failedLinks)

	for err := range failedLinks {
		fmt.Println(err)
	}

	return nil
}
