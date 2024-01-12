package cmd

import (
	"fmt"
	"os"

	"github.com/michalczmiel/batch-image-getter/internal"
	"github.com/spf13/cobra"
)

var htmlCmd = &cobra.Command{
	Use:   "html <url>",
	Short: "Download all images from an HTML website",
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
	htmlCmd.Flags().StringArrayP("types", "t", []string{".jpg", ".jpeg", ".png"}, "image types to download")
	htmlCmd.Flags().IntP("concurrency", "c", 10, "number of concurrent downloads")
	htmlCmd.Flags().StringP("dir", "d", internal.DefaultPath, "directory to save images to")
	htmlCmd.Flags().String("user-agent", "", "custom user agent to use for requests")
	rootCmd.AddCommand(htmlCmd)
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

	directory, err := cmd.Flags().GetString("dir")
	if err != nil {
		return err
	}

	userAgent, err := cmd.Flags().GetString("user-agent")
	if err != nil {
		return err
	}

	err = internal.DownloadImagesFromWebsite(url, internal.Parameters{
		Directory:  directory,
		ImageTypes: imageTypesToDownload,
		Concurrent: concurrentWorkersCount,
		UserAgent:  userAgent,
	})
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
