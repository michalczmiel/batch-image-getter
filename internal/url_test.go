package internal

import (
	"testing"
)

func TestGetFileNameFromUrl(t *testing.T) {
	testdata := []struct {
		url      string
		expected string
	}{
		{"https://example.com/logo_main-v2.png", "logo_main-v2.png"},
		{"https://example.com/resource-1240-720", "resource-1240-720"},
		{"https://example.com/image.jpg?w=1919&h=1280", "image.jpg"},
	}

	for _, tt := range testdata {
		actual, err := GetFileNameFromUrl(tt.url)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if actual != tt.expected {
			t.Errorf("want %s got %s", tt.expected, actual)
		}
	}
}

func TestProcessLinks(t *testing.T) {
	rawLinks := []string{
		"/logo_main-v2.png",
		"/2023/01/30/cover.png",
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
		"https://example.com/logo.svg",
		"https://example.com/image3.jpg",
		"https://example.com/image4.JPG",
		"https://example.com/image5.Png",
	}

	actual := ProcessLinks(url, rawLinks)

	if len(actual) != len(expected) {
		t.Errorf("want %v got %v", expected, actual)
	}
}

func TestRemoveDuplicates(t *testing.T) {
	given := []string{"a", "b", "b", "c", "d", "a", "b", "e"}
	expected := []string{"a", "b", "c", "d", "e"}

	actual := RemoveDuplicates(given)

	if len(expected) != len(actual) {
		t.Errorf("want %v got %v", expected, actual)
	}
}

func TestGetRootUrl(t *testing.T) {
	testdata := []struct {
		url      string
		expected string
	}{
		{"https://example.com", "https://example.com"},
		{"http://example.com", "http://example.com"},
		{"https://example.com/images/image.jpg", "https://example.com"},
		{"https://example.com/images/image2.jpg?w=1919&h=1280", "https://example.com"},
		{"https://example.com?example=query", "https://example.com"},
	}

	for _, tt := range testdata {
		actual := getRootUrl(tt.url)

		if actual != tt.expected {
			t.Errorf("want %s got %s", tt.expected, actual)
		}
	}
}
