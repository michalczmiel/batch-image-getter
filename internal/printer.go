package internal

import (
	"fmt"
	"io"
	"os"
)

type Printer struct {
	writer io.Writer
}

func NewStdoutPrinter() *Printer {
	return &Printer{writer: os.Stdout}
}

func (p *Printer) PrintResultsAsPlainText(results []DownloadResult) {
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
