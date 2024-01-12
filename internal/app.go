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

type Parameters struct {
	Directory  string
	ImageTypes []string
	Concurrent int
	UserAgent  string
}

func downloadWorker(wg *sync.WaitGroup, parameters Parameters, linksToProcess <-chan string, failedLinks chan<- error) {
	defer wg.Done()

	for link := range linksToProcess {
		fileName := GetFileNameFromUrl(link)
		fmt.Printf("Downloading %s\n", link)

		filePath := path.Join(parameters.Directory, fileName)

		err := DownloadFileFromUrl(link, filePath, parameters.UserAgent)

		if err != nil {
			failedLinks <- fmt.Errorf("error downloading file %s %v", link, err)
		}
	}
}

func DownloadImagesFromWebsite(url string, parameters Parameters) error {
	doc, err := GetHtmlDocFromUrl(url, parameters.UserAgent)
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
		go downloadWorker(&wg, parameters, linksToProcess, failedLinks)
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
