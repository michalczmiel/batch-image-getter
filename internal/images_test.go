package internal

import (
	"fmt"
	"testing"
)

func TestPrepareLinksForDownload(t *testing.T) {
	urls := []string{
		"https://example.com/v2/file/5a9abk2-PL/image;s=200x0;q=50",
		"https://example.com/v2/file/fad6g2x-PL/image;s=200x0;q=50",
	}

	parameters := &Parameters{
		Directory: "images",
	}

	given := PrepareLinksForDownload(urls, parameters)

	expected := []DownloadInput{
		{Url: "https://example.com/v2/file/5a9abk2-PL/image;s=200x0;q=50", FilePath: "images/image;s=200x0;q=50"},
		{Url: "https://example.com/v2/file/fad6g2x-PL/image;s=200x0;q=50", FilePath: "images/1image;s=200x0;q=50"},
	}

	for i, input := range given {
		if input.Url != expected[i].Url {
			t.Errorf("got %v expected %v", input.Url, expected[i].Url)
		}

		if input.FilePath != expected[i].FilePath {
			t.Errorf("got %v expected %v", input.FilePath, expected[i].FilePath)
		}
	}
}

func TestGetImageType(t *testing.T) {
	testdata := []struct {
		contentType string
		expected    string
		error       error
	}{
		{"image/jpeg", "jpeg", nil},
		{"image/png", "png", nil},
		{"image/webp", "webp", nil},
		{"", "", fmt.Errorf("content type is empty")},
		{"text/html", "", fmt.Errorf("content type 'text/html' is not an image")},
	}

	for _, data := range testdata {
		imageType, err := getImageType(data.contentType)

		if err != nil && data.error == nil {
			t.Errorf("got %v expected %v", err, data.error)
		}

		if imageType != data.expected {
			t.Errorf("got %v expected %v", imageType, data.expected)
		}
	}
}

func TestAddExtensionIfMissing(t *testing.T) {
	testdata := []struct {
		filePath  string
		imageType string
		expected  string
	}{
		{"image", "jpg", "image.jpg"},
		{"image.jpg", "jpg", "image.jpg"},
		{"image.png", "png", "image.png"},
		{"image.webp", "webp", "image.webp"},
	}

	for _, data := range testdata {
		given := addExtensionIfMissing(data.filePath, data.imageType)

		if given != data.expected {
			t.Errorf("got %v expected %v", given, data.expected)
		}
	}
}
