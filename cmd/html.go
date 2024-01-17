package cmd

import (
	"fmt"

	"github.com/michalczmiel/batch-image-getter/internal"
	"github.com/spf13/cobra"
)

var htmlCmd = &cobra.Command{
	Use:   "html <url>",
	Short: "Download all images from an HTML website",
	RunE:  runHtmlCmd,
	Args: func(cmd *cobra.Command, args []string) error {
		err := validateArguments(args)
		if err != nil {
			return err
		}

		concurrentWorkersCount, err := cmd.Flags().GetInt("concurrency")
		if err != nil {
			return err
		}

		if concurrentWorkersCount < 1 {
			return fmt.Errorf("concurrency must be greater than 0")
		}

		return nil
	},
}

func runHtmlCmd(cmd *cobra.Command, args []string) error {
	url := args[0]

	parameters, err := getRootParameters(cmd)
	if err != nil {
		return err
	}

	doc, err := internal.GetHtmlDocFromUrl(url, parameters.UserAgent)
	if err != nil {
		return err
	}

	rawLinks := internal.GetImageLinksFromHtmlDoc(doc)
	if len(rawLinks) == 0 {
		return fmt.Errorf("no links found")
	}

	links := internal.ProcessLinks(url, rawLinks)

	fmt.Printf("Found %d valid image links\n", len(links))

	err = internal.CreateDirectoryIfDoesNotExists(parameters.Directory)
	if err != nil {
		return err
	}

	err = internal.DownloadImages(links, parameters)
	if err != nil {
		return err
	}

	return nil

}

func validateArguments(args []string) error {
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
	rootCmd.AddCommand(htmlCmd)
}
