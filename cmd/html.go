package cmd

import (
	"fmt"

	"github.com/michalczmiel/batch-image-getter/internal"
	"github.com/spf13/cobra"
)

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

var htmlCmd = &cobra.Command{
	Use:   "html <url>",
	Short: "Download all images from an HTML website",
	RunE:  runHtmlCmd,
	Args: func(cmd *cobra.Command, args []string) error {
		err := validateArguments(args)
		if err != nil {
			return err
		}

		err = validateRootFlags(cmd)
		if err != nil {
			return err
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

	httClient := internal.NewHttpClient(parameters.UserAgent)

	doc, err := internal.GetHtmlDocFromUrl(url, httClient, parameters)
	if err != nil {
		return err
	}

	rawLinks := internal.GetImageLinksFromHtmlDoc(doc)
	if len(rawLinks) == 0 {
		return fmt.Errorf("no links found")
	}

	links := internal.ProcessLinks(url, rawLinks)
	links = internal.RemoveDuplicates(links)

	printer := internal.NewStdoutPrinter(parameters.OutputFormat)
	printer.PrintProgress(len(links))

	err = internal.CreateDirectoryIfDoesNotExists(parameters.Directory)
	if err != nil {
		return err
	}

	inputs := internal.PrepareLinksForDownload(links, parameters)
	results := internal.DownloadImages(inputs, httClient, parameters)
	printer.PrintResults(results)

	return nil
}

func init() {
	rootCmd.AddCommand(htmlCmd)
}
