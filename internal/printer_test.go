package internal

import (
	"bytes"
	"fmt"
	"testing"
)

func Results() []DownloadResult {
	return []DownloadResult{
		{Url: "https://example.com/image1.jpg", Path: "images/image1.jpg", Err: nil},
		{Url: "https://example.com/image2.jpg", Path: "", Err: fmt.Errorf("invalid content type")},
		{Url: "https://example.com/image3.jpg", Path: "", Err: fmt.Errorf("error downloading image")},
		{Url: "https://example.com/image4.jpg", Path: "images/image1.jpg", Err: nil},
	}
}

func TestPlainTextPrinter(t *testing.T) {
	t.Run("Prints results", func(t *testing.T) {
		var buf bytes.Buffer

		printer := PlainTextPrinter{writer: &buf}

		printer.PrintResults(Results())

		got := buf.String()

		expected := "Successful downloads: 2\nFailed downloads: 2\n"

		if got != expected {
			t.Errorf("got %q expected %q", got, expected)
		}
	})

	t.Run("Prints progress", func(t *testing.T) {
		var buf bytes.Buffer

		printer := PlainTextPrinter{writer: &buf}

		printer.PrintProgress(15)

		got := buf.String()

		expected := "Found 15 valid image links\n"

		if got != expected {
			t.Errorf("got %q expected %q", got, expected)
		}
	})
}

func TestJsonPrinter(t *testing.T) {
	t.Run("Prints results", func(t *testing.T) {
		var buf bytes.Buffer

		printer := JsonPrinter{writer: &buf}

		printer.PrintResults(Results())

		got := buf.String()

		expected := "[{\"url\":\"https://example.com/image1.jpg\",\"path\":\"images/image1.jpg\"},{\"url\":\"https://example.com/image2.jpg\",\"error\":\"invalid content type\"},{\"url\":\"https://example.com/image3.jpg\",\"error\":\"error downloading image\"},{\"url\":\"https://example.com/image4.jpg\",\"path\":\"images/image1.jpg\"}]"

		if got != expected {
			t.Errorf("got %q expected %q", got, expected)
		}
	})
}
