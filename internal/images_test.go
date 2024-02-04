package internal

import "testing"

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
