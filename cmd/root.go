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

func getRootParameters(cmd *cobra.Command) (*internal.Parameters, error) {
	imageTypesToDownload, err := cmd.Flags().GetStringArray("types")
	if err != nil {
		return nil, err
	}

	concurrentWorkersCount, err := cmd.Flags().GetInt("concurrency")
	if err != nil {
		return nil, err
	}

	directory, err := cmd.Flags().GetString("dir")
	if err != nil {
		return nil, err
	}

	userAgent, err := cmd.Flags().GetString("user-agent")
	if err != nil {
		return nil, err
	}

	referer, err := cmd.Flags().GetString("referer")
	if err != nil {
		return nil, err
	}

	parameters := &internal.Parameters{
		ImageTypes: imageTypesToDownload,
		Directory:  directory,
		Concurrent: concurrentWorkersCount,
		UserAgent:  userAgent,
		Referer:    referer,
	}

	return parameters, nil
}

func validateRootFlags(cmd *cobra.Command) error {
	concurrentWorkersCount, err := cmd.Flags().GetInt("concurrency")
	if err != nil {
		return err
	}

	if concurrentWorkersCount < 1 {
		return fmt.Errorf("concurrency must be greater than 0")
	}

	return nil
}

func addRootFlags(cmd *cobra.Command) {
	cmd.Flags().StringArrayP("types", "t", []string{"jpg", "jpeg", "png", "gif", "webp"}, "image types to download")
	cmd.Flags().IntP("concurrency", "c", 10, "number of concurrent downloads")
	cmd.Flags().StringP("dir", "d", internal.DefaultPath, "directory to save images to")
	cmd.Flags().String("user-agent", "", "custom user agent to use for requests")
	cmd.Flags().String("referer", "", "custom referer to use for requests")
}

func init() {
	addRootFlags(htmlCmd)
	addRootFlags(fileCmd)
}

func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
