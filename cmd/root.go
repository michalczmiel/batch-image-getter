package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/michalczmiel/batch-image-downloader/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	RunE: run,
	Args: cobra.MinimumNArgs(1),
}

func init() {
	rootCmd.Flags().StringArrayP("types", "t", []string{".jpg", ".jpeg", ".png"}, "image types to download")
}

func run(cmd *cobra.Command, args []string) error {
	url := args[0]

	imageTypesToDownload, err := cmd.Flags().GetStringArray("types")
	if err != nil {
		return err
	}

	doc, err := internal.GetHtmlDocFromUrl(url)
	if err != nil {
		return err
	}

	links := internal.GetImageLinksFromHtmlDoc(doc, imageTypesToDownload)
	if len(links) == 0 {
		return fmt.Errorf("no links found")
	}

	var wg sync.WaitGroup
	for _, link := range links {
		wg.Add(1)

		go func(l string) {
			defer wg.Done()

			fileName := internal.GetFileNameFromUrl(l)
			err = internal.DownloadFileFromUrl(l, fileName)
			if err != nil {
				fmt.Printf("Error downloading file %s %v", l, err)
			}
		}(link)
	}

	wg.Wait()

	return nil
}

func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
