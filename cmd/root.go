package cmd

import (
	"fmt"
	"os"

	"github.com/michalczmiel/batch-image-getter/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "big",
	Short: "big is a CLI tool for downloading images",
}

func getRootParameters(cmd *cobra.Command) (internal.Parameters, error) {
	imageTypesToDownload, err := cmd.Flags().GetStringArray("types")
	if err != nil {
		return internal.Parameters{}, err
	}

	concurrentWorkersCount, err := cmd.Flags().GetInt("concurrency")
	if err != nil {
		return internal.Parameters{}, err
	}

	directory, err := cmd.Flags().GetString("dir")
	if err != nil {
		return internal.Parameters{}, err
	}

	userAgent, err := cmd.Flags().GetString("user-agent")
	if err != nil {
		return internal.Parameters{}, err
	}

	parameters := internal.Parameters{
		ImageTypes: imageTypesToDownload,
		Directory:  directory,
		Concurrent: concurrentWorkersCount,
		UserAgent:  userAgent,
	}

	return parameters, nil
}

func init() {
	htmlCmd.Flags().StringArrayP("types", "t", []string{"jpg", "jpeg", "png", "gif", "webp"}, "image types to download")
	htmlCmd.Flags().IntP("concurrency", "c", 10, "number of concurrent downloads")
	htmlCmd.Flags().StringP("dir", "d", internal.DefaultPath, "directory to save images to")
	htmlCmd.Flags().String("user-agent", "", "custom user agent to use for requests")
}

func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
