package internal

import "testing"

func TestGetFileNameFromUrl(t *testing.T) {
	url := "https://example.com/2023/01/01/logo_main-v2.png"
	expected := "logo_main-v2.png"
	actual := GetFileNameFromUrl(url)

	if actual != expected {
		t.Errorf("want %s got %s", actual, expected)
	}
}
