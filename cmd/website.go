package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/michalczmiel/batch-image-downloader/internal"
	"github.com/spf13/cobra"
)

var websiteCmd = &cobra.Command{
	Use:   "website <url>",
	Short: "Download all images from a website",
	RunE:  run,
	Args:  validateWebsiteCmdArgs,
}

func validateWebsiteCmdArgs(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("requires a url argument")
	}

	if len(args) > 1 {
		return fmt.Errorf("too many arguments, please provide a single url")
	}

	url := args[0]
	if !internal.IsUrlValid(url) {
		return fmt.Errorf("invalid url")
	}

	return nil
}

func init() {
	websiteCmd.Flags().StringArrayP("types", "t", []string{".jpg", ".jpeg", ".png"}, "image types to download")
	rootCmd.AddCommand(websiteCmd)
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

	rawLinks := internal.GetImageLinksFromHtmlDoc(doc)
	if len(rawLinks) == 0 {
		return fmt.Errorf("no links found")
	}

	links := internal.ProcessLinks(url, rawLinks, imageTypesToDownload)

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
