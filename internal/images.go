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

type DownloadInput struct {
	Url      string
	FilePath string
}

func PrepareLinksForDownload(links []string, parameters *Parameters) []DownloadInput {
	var downloadInputs []DownloadInput

	for _, link := range links {
		fileName, err := GetFileNameFromUrl(link)
		if err != nil {
			continue
		}

		filePath := path.Join(parameters.Directory, fileName)
		downloadInputs = append(downloadInputs, DownloadInput{Url: link, FilePath: filePath})
	}

	return downloadInputs
}

func downloadImage(link DownloadInput, httClient HttpClent, parameters *Parameters) error {
	var referer string
	if parameters.Referer == "" {
		referer = getRootUrl(link.Url)
	} else {
		referer = parameters.Referer
	}

	response, err := httClient.Request(link.Url, map[string]string{
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

	filePath := addExtensionIfMissing(link.FilePath, contentType)

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

func DownloadImages(links []DownloadInput, httClient HttpClent, parameters *Parameters) []DownloadResult {
	linksToProcess := make(chan DownloadInput, len(links))
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
					results <- DownloadResult{Url: link.Url, Err: err}
				} else {
					results <- DownloadResult{Url: link.Url, Err: nil}
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
