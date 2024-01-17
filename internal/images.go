package internal

import (
	"fmt"
	"path"
	"sync"
)

type Parameters struct {
	ImageTypes []string
	Directory  string
	Concurrent int
	UserAgent  string
}

func downloadWorker(wg *sync.WaitGroup, parameters Parameters, linksToProcess <-chan string, failedLinks chan<- error) {
	for link := range linksToProcess {
		fileName := GetFileNameFromUrl(link)

		filePath := path.Join(parameters.Directory, fileName)

		err := DownloadImageFromUrl(link, filePath, parameters)

		if err != nil {
			failedLinks <- fmt.Errorf("error downloading file %s %v", link, err)
		}

		wg.Done()
	}
}

func DownloadImages(links []string, parameters Parameters) error {
	linksToProcess := make(chan string, len(links))
	failedLinks := make(chan error, len(links))
	var wg sync.WaitGroup

	for _, link := range links {
		linksToProcess <- link
	}
	close(linksToProcess)

	// all links need to be processed before closing
	wg.Add(len(links))

	for i := 0; i < parameters.Concurrent; i++ {
		go downloadWorker(&wg, parameters, linksToProcess, failedLinks)
	}

	// wait for all workers to finish
	go func() {
		wg.Wait()
		close(failedLinks)
	}()

	for err := range failedLinks {
		fmt.Println(err)
	}

	return nil
}
