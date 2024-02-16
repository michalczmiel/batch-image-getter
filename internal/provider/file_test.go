package provider

import (
	"io"
	"testing"
)

type mockFileSystem struct{}

func (m mockFileSystem) Exists(path string) bool {
	return path == "images.txt"
}

func (m mockFileSystem) ReadLines(path string) ([]string, error) {
	return []string{
		"https://example.com/image1.jpg",
		"-",
		"\n",
		"https://example.com/image2.jpg",
		"https://example.com/image3.jpg",
	}, nil
}

func (m mockFileSystem) WriteLines(path string, lines []string) error {
	return nil
}

func (m mockFileSystem) CreateDirectory(path string) error {
	return nil
}

func (m mockFileSystem) Save(body io.ReadCloser, path string) error {
	return nil
}

func TestProviderLinks(t *testing.T) {
	t.Run("should return error when file does not exist", func(t *testing.T) {
		fileSystem := mockFileSystem{}

		provider := NewFileProvider("non-existing-file.txt", fileSystem)

		_, err := provider.Links()

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("should return valid links from existing file", func(t *testing.T) {
		fileSystem := mockFileSystem{}

		provider := NewFileProvider("images.txt", fileSystem)

		links, err := provider.Links()

		if err != nil {
			t.Error(err)
		}

		expectedLinks := []string{
			"https://example.com/image1.jpg",
			"https://example.com/image2.jpg",
			"https://example.com/image3.jpg",
		}

		if len(links) != len(expectedLinks) {
			t.Fatalf("Expected %d links, got %d", len(expectedLinks), len(links))
		}

		for i, link := range links {
			if link != expectedLinks[i] {
				t.Fatalf("Expected %s, got %s", expectedLinks[i], link)
			}
		}
	})
}
