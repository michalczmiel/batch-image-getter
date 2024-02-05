package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type OutputFormat string

const (
	PlainText OutputFormat = "plain"
	Json      OutputFormat = "json"
)

type Printer interface {
	PrintResults(results []DownloadResult)
	PrintProgress(validLinks int)
}

type JsonPrinter struct {
	writer io.Writer
}

type resultJson struct {
	Url   string `json:"url"`
	Error string `json:"error,omitempty"`
	Path  string `json:"path,omitempty"`
}

func (p *JsonPrinter) PrintResults(results []DownloadResult) {
	resultsOutput := make([]resultJson, 0, len(results))

	for _, result := range results {
		var output resultJson

		if result.Err != nil {
			output = resultJson{Url: result.Url, Error: result.Err.Error()}
		} else {
			output = resultJson{Url: result.Url, Path: result.Path}
		}

		resultsOutput = append(resultsOutput, output)
	}

	jsonData, err := json.Marshal(resultsOutput)
	if err != nil {
		fmt.Fprintf(p.writer, "unexpected error: %s", err)
		return
	}

	fmt.Fprintf(p.writer, "%s", jsonData)
}

func (p *JsonPrinter) PrintProgress(validLinks int) {
	// no-op
}

type PlainTextPrinter struct {
	writer io.Writer
}

func (p *PlainTextPrinter) PrintResults(results []DownloadResult) {
	successfulDownloads := 0
	failedDownloads := 0

	for _, result := range results {
		if result.Err != nil {
			failedDownloads++
		} else {
			successfulDownloads++
		}
	}

	fmt.Fprintf(p.writer, "Successful downloads: %d\n", successfulDownloads)
	fmt.Fprintf(p.writer, "Failed downloads: %d\n", failedDownloads)
}

func (p *PlainTextPrinter) PrintProgress(validLinks int) {
	fmt.Fprintf(p.writer, "Found %d valid image links\n", validLinks)
}

func NewStdoutPrinter(outputFormat OutputFormat) Printer {
	if outputFormat == Json {
		return &JsonPrinter{writer: os.Stdout}
	}

	return &PlainTextPrinter{writer: os.Stdout}
}
