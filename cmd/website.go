package cmd

import (
	"fmt"
	"os"

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

	err = internal.DownloadImagesFromWebsite(url, imageTypesToDownload)
	if err != nil {
		return err
	}

	return nil
}

func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
