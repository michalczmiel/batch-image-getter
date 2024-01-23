package internal

import (
	"path"
	"sync"
)

type Parameters struct {
	ImageTypes []string
	Directory  string
	Concurrent int
	UserAgent  string
}

func downloadWorker(wg *sync.WaitGroup, parameters *Parameters, linksToProcess <-chan string, results chan<- DownloadResult) {
	for link := range linksToProcess {
		fileName, err := GetFileNameFromUrl(link)
		if err != nil {
			results <- DownloadResult{Url: link, Err: err}
		}

		filePath := path.Join(parameters.Directory, fileName)

		err = DownloadImageFromUrl(link, filePath, parameters)

		if err != nil {
			results <- DownloadResult{Url: link, Err: err}
		}

		results <- DownloadResult{Url: link, Err: nil}

		wg.Done()
	}
}

type DownloadResult struct {
	Url string
	Err error
}

func DownloadImages(links []string, parameters *Parameters) []DownloadResult {
	linksToProcess := make(chan string, len(links))
	results := make(chan DownloadResult, len(links))

	var wg sync.WaitGroup

	for _, link := range links {
		linksToProcess <- link
	}
	close(linksToProcess)

	// all links need to be processed before closing
	wg.Add(len(links))

	for i := 0; i < parameters.Concurrent; i++ {
		go downloadWorker(&wg, parameters, linksToProcess, results)
	}

	// wait for all workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	resultsOutput := []DownloadResult{}
	for result := range results {
		resultsOutput = append(resultsOutput, result)
	}
	return resultsOutput
}
