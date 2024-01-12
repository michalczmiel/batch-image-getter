package internal

import (
	"fmt"
	"os"
	"path"
	"sync"
)

const DefaultPath = "."

func createDirectoryIfDoesNotExists(directory string) error {
	if directory == DefaultPath {
		return nil
	}

	if directory == "" {
		directory = DefaultPath
	}

	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directory %s %v", directory, err)
	}

	return nil
}

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

type Parameters struct {
	Directory  string
	ImageTypes []string
	Concurrent int
}

func DownloadImagesFromWebsite(url string, parameters Parameters) error {
	doc, err := GetHtmlDocFromUrl(url)
	if err != nil {
		return err
	}

	rawLinks := GetImageLinksFromHtmlDoc(doc)
	if len(rawLinks) == 0 {
		return fmt.Errorf("no links found")
	}

	links := ProcessLinks(url, rawLinks, parameters.ImageTypes)

	fmt.Printf("Found %d valid image links\n", len(links))

	err = createDirectoryIfDoesNotExists(parameters.Directory)
	if err != nil {
		return err
	}

	linksToProcess := make(chan string)
	failedLinks := make(chan error, parameters.Concurrent)

	var wg sync.WaitGroup

	for i := 0; i < parameters.Concurrent; i++ {
		wg.Add(1)
		go downloadWorker(&wg, parameters.Directory, linksToProcess, failedLinks)
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
