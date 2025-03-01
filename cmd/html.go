package cmd

import (
	"fmt"

	"github.com/michalczmiel/batch-image-getter/internal"
	"github.com/michalczmiel/batch-image-getter/internal/provider"
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

	regexSearch, err := cmd.Flags().GetBool("regex-search")
	if err != nil {
		return err
	}

	fileSystem := internal.NewFileSystem()
	httpClient := internal.NewHttpClient(parameters.UserAgent)

	provider := provider.NewHtmlProvider(url, httpClient, parameters, regexSearch)

	links, err := provider.Links()
	if err != nil {
		return err
	}

	referer, err := internal.GetRootUrl(url)
	if err != nil {
		return err
	}

	// download images from the same domain
	parameters.Referer = referer

	return internal.Run(links, parameters, httpClient, fileSystem)
}

func init() {
	htmlCmd.Flags().Bool("regex-search", false, "use regex to search for images in the HTML content")

	rootCmd.AddCommand(htmlCmd)
}
