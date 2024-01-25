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
	Referer    string
}

func downloadImage(link string, httClient HttpClent, parameters *Parameters) error {
	fileName, err := GetFileNameFromUrl(link)
	if err != nil {
		return err
	}

	filePath := path.Join(parameters.Directory, fileName)

	var referer string
	if parameters.Referer == "" {
		referer = getRootUrl(link)
	} else {
		referer = parameters.Referer
	}

	response, err := httClient.Request(link, map[string]string{
		"User-Agent": parameters.UserAgent,
		"Referer":    referer,
	})
	if err != nil {
		return err
	}

	contentType := response.Header.Get("Content-Type")
	err = validateContentType(contentType, parameters.ImageTypes)
	if err != nil {
		return err
	}

	filePath = addExtensionIfMissing(filePath, contentType)

	err = SaveToFile(response.Body, filePath)
	if err != nil {
		return err
	}

	return nil
}

type DownloadResult struct {
	Url string
	Err error
}

func DownloadImages(links []string, httClient HttpClent, parameters *Parameters) []DownloadResult {
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
		go func() {
			for link := range linksToProcess {
				err := downloadImage(link, httClient, parameters)

				if err != nil {
					results <- DownloadResult{Url: link, Err: err}
				} else {
					results <- DownloadResult{Url: link, Err: nil}
				}

				wg.Done()
			}
		}()
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
