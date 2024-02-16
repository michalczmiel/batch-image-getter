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

	json, err := cmd.Flags().GetBool("json")
	if err != nil {
		return nil, err
	}

	sleepInterval, err := cmd.Flags().GetInt("sleep-interval")
	if err != nil {
		return nil, err
	}

	maxSleepInterval, err := cmd.Flags().GetInt("max-sleep-interval")
	if err != nil {
		return nil, err
	}

	outputFormat := internal.PlainText
	if json {
		outputFormat = internal.Json
	}

	parameters := &internal.Parameters{
		ImageTypes:       imageTypesToDownload,
		Directory:        directory,
		Concurrent:       concurrentWorkersCount,
		UserAgent:        userAgent,
		Referer:          referer,
		OutputFormat:     outputFormat,
		SleepInterval:    sleepInterval,
		MaxSleepInterval: maxSleepInterval,
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

	sleepInterval, err := cmd.Flags().GetInt("sleep-interval")
	if err != nil {
		return err
	}

	if sleepInterval < 0 {
		return fmt.Errorf("sleep interval must be greater than or equal to 0")
	}

	maxSleepInterval, err := cmd.Flags().GetInt("max-sleep-interval")
	if err != nil {
		return err
	}

	if maxSleepInterval < 0 {
		return fmt.Errorf("max sleep interval must be greater than or equal to 0")
	}

	if maxSleepInterval > 0 && maxSleepInterval < sleepInterval {
		return fmt.Errorf("max sleep interval must be greater than or equal to sleep interval")
	}

	return nil
}

func addRootFlags(cmd *cobra.Command) {
	cmd.Flags().StringArray("types", []string{"jpg", "jpeg", "png", "gif", "webp"}, "image types to download")
	cmd.Flags().IntP("concurrency", "c", 10, "number of concurrent downloads")
	cmd.Flags().StringP("dir", "d", internal.DefaultPath, "directory to save images to")
	cmd.Flags().String("user-agent", "", "custom user agent to use for requests")
	cmd.Flags().String("referer", "", "custom referer to use for requests")
	cmd.Flags().Bool("json", false, "output results as json")
	cmd.Flags().Int("sleep-interval", 0, "number of seconds to sleep after each request or minimum number of seconds for randomized sleep when used along max-sleep-interval (default 0)")
	cmd.Flags().Int("max-sleep-interval", 0, "maximum number of seconds to sleep after each request")
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
