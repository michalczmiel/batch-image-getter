package internal

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPrintResultsAsPlainText(t *testing.T) {
	var buf bytes.Buffer

	printer := Printer{writer: &buf}

	results := []DownloadResult{
		{Url: "https://example.com/image1.jpg", Err: nil},
		{Url: "https://example.com/image2.jpg", Err: fmt.Errorf("invalid content type")},
		{Url: "https://example.com/image3.jpg", Err: fmt.Errorf("error downloading image")},
		{Url: "https://example.com/image4.jpg", Err: nil},
	}

	printer.PrintResultsAsPlainText(results)

	got := buf.String()

	expected := "Successful downloads: 2\nFailed downloads: 2\n"

	if got != expected {
		t.Errorf("got %q expected %q", got, expected)
	}
}
