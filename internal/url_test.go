package internal

import (
	"testing"
)

func TestGetFileNameFromUrl(t *testing.T) {
	url := "https://example.com/2023/01/01/logo_main-v2.png"
	expected := "logo_main-v2.png"
	actual := GetFileNameFromUrl(url)

	if actual != expected {
		t.Errorf("want %s got %s", actual, expected)
	}
}

func TestProcessLinks(t *testing.T) {
	fileTypes := []string{".jpg", ".jpeg", ".png"}

	rawLinks := []string{
		"/logo_main-v2.png",
		"/2023/01/30/cover.png",
		"https://example.com/image1.jpg",
		"https://example.com/image1.jpg",
		"http://example.com/image2.jpg",
		"//example.com/image3.jpg",
		"https://example.com/logo.svg",
		"https://example.com/image4.JPG",
		"https://example.com/image5.Png",
	}
	url := "https://example.com"
	expected := []string{
		"https://example.com/logo_main-v2.png",
		"https://example.com/2023/01/30/cover.png",
		"https://example.com/image1.jpg",
		"http://example.com/image2.jpg",
		"https://example.com/image3.jpg",
		"https://example.com/image4.JPG",
		"https://example.com/image5.Png",
	}

	actual := ProcessLinks(url, rawLinks, fileTypes)

	if len(actual) != len(expected) {
		t.Errorf("want %v got %v", expected, actual)
	}
}
