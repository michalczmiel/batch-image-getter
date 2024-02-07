package internal

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

type Parameters struct {
	ImageTypes   []string
	Directory    string
	Concurrent   int
	UserAgent    string
	Referer      string
	OutputFormat OutputFormat
}

type DownloadInput struct {
	Url      string
	FilePath string
}

func PrepareLinksForDownload(links []string, parameters *Parameters) []DownloadInput {
	var downloadInputs []DownloadInput

	// set is not available in Go, so we use map instead to remove duplicates
	var alreadyTakenNames = map[string]struct{}{}

	for index, link := range links {
		fileName, err := GetFileNameFromUrl(link)
		if err != nil {
			continue
		}

		_, exists := alreadyTakenNames[fileName]
		if exists {
			fileName = fmt.Sprint(index) + fileName
		}

		alreadyTakenNames[fileName] = struct{}{}

		filePath := path.Join(parameters.Directory, fileName)
		downloadInputs = append(downloadInputs, DownloadInput{Url: link, FilePath: filePath})
	}

	return downloadInputs
}

func getImageType(contentType string) (string, error) {
	if contentType == "" {
		return "", fmt.Errorf("content type is empty")
	}

	if !strings.HasPrefix(contentType, "image") {
		return "", fmt.Errorf("content type '%s' is not an image", contentType)
	}

	imageType := strings.Split(contentType, "/")[1]

	return imageType, nil
}

func isImageTypeAllowed(imageType string, allowedImageTypes []string) bool {
	for _, allowedImageType := range allowedImageTypes {
		if imageType == allowedImageType {
			return true
		}
	}

	return false
}

func addExtensionIfMissing(filePath, imageType string) string {
	extension := filepath.Ext(filePath)

	if extension != "" {
		return filePath
	}

	extension = "." + imageType

	return filePath + extension
}

func downloadImage(link DownloadInput, httClient HttpClient, fileSystem FileSystem, parameters *Parameters) (outputPath string, err error) {
	referer, err := getRootUrl(link.Url)
	if err != nil {
		return "", err
	}

	if parameters.Referer != "" {
		referer = parameters.Referer
	}

	response, err := httClient.Request(link.Url, map[string]string{
		"Referer": referer,
	})
	if err != nil {
		return "", err
	}

	imageType, err := getImageType(response.Header.Get("Content-Type"))
	if err != nil {
		return "", err
	}

	if !isImageTypeAllowed(imageType, parameters.ImageTypes) {
		return "", fmt.Errorf("image type '%s' is not allowed", imageType)
	}

	filePath := addExtensionIfMissing(link.FilePath, imageType)

	err = fileSystem.Save(response.Body, filePath)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

type DownloadResult struct {
	Url  string
	Err  error
	Path string
}

func DownloadImages(links []DownloadInput, httClient HttpClient, fileSystem FileSystem, parameters *Parameters) []DownloadResult {
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
				outputPath, err := downloadImage(link, httClient, fileSystem, parameters)

				if err != nil {
					results <- DownloadResult{Url: link.Url, Err: err, Path: ""}
				} else {
					results <- DownloadResult{Url: link.Url, Err: nil, Path: outputPath}
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
