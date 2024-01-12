package internal

import (
	"fmt"
	"path"
	"sync"
)

func downloadWorker(wg *sync.WaitGroup, directory string, linksToProcess <-chan string, failedLinks chan<- error) {
	defer wg.Done()

	for link := range linksToProcess {
		fileName := GetFileNameFromUrl(link)
		fmt.Printf("Downloading %s\n", link)

		filePath := path.Join(directory, fileName)

		err := DownloadFileFromUrl(link, filePath)

		if err != nil {
			failedLinks <- fmt.Errorf("error downloading file %s %v", link, err)
		}
	}
}

// TODO: refactor arguments to be struct
func DownloadImagesFromWebsite(url string, imageTypesToDownload []string, concurrentWorkersCount int, directory string) error {
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

	linksToProcess := make(chan string)
	failedLinks := make(chan error, concurrentWorkersCount)

	var wg sync.WaitGroup

	for i := 0; i < concurrentWorkersCount; i++ {
		wg.Add(1)
		go downloadWorker(&wg, directory, linksToProcess, failedLinks)
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
