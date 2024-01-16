package internal

import (
	"fmt"
	"path"
	"sync"
)

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

		err := DownloadImageFromUrl(link, filePath, parameters.UserAgent)

		if err != nil {
			failedLinks <- fmt.Errorf("error downloading file %s %v", link, err)
		}
	}
}

func DownloadImages(links []string, parameters Parameters) error {
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
