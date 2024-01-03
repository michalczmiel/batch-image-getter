package cmd

import (
	"fmt"
	"os"

	"github.com/michalczmiel/batch-image-getter/internal"
	"github.com/spf13/cobra"
)

var websiteCmd = &cobra.Command{
	Use:   "website <url>",
	Short: "Download all images from a website",
	RunE:  run,
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
	websiteCmd.Flags().StringArrayP("types", "t", []string{".jpg", ".jpeg", ".png"}, "image types to download")
	websiteCmd.Flags().IntP("concurrency", "c", 10, "number of concurrent downloads")
	rootCmd.AddCommand(websiteCmd)
}

func run(cmd *cobra.Command, args []string) error {
	url := args[0]

	imageTypesToDownload, err := cmd.Flags().GetStringArray("types")
	if err != nil {
		return err
	}

	concurrentWorkersCount, err := cmd.Flags().GetInt("concurrency")
	if err != nil {
		return err
	}

	err = internal.DownloadImagesFromWebsite(url, imageTypesToDownload, concurrentWorkersCount)
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
