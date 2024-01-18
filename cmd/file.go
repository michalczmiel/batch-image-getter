package cmd

import (
	"fmt"
	"os"

	"github.com/michalczmiel/batch-image-getter/internal"
	"github.com/spf13/cobra"
)

func validateFileCmdArguments(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("requires a path argument")
	}

	if len(args) > 1 {
		return fmt.Errorf("too many arguments, please provide a single file path")
	}

	path := args[0]
	_, err := os.Stat(path)

	if err != nil {
		return fmt.Errorf("file does not exist")
	}

	return nil
}

var fileCmd = &cobra.Command{
	Use:   "file <path>",
	Short: "Download all imagess from a file",
	RunE:  runFileCmd,
	Args: func(cmd *cobra.Command, args []string) error {
		err := validateFileCmdArguments(args)
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

func runFileCmd(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	parameters, err := getRootParameters(cmd)
	if err != nil {
		return err
	}

	links, err := internal.GetLinesFromFile(filePath)
	if err != nil {
		return err
	}

	links = internal.RemoveDuplicates(links)

	if len(links) == 0 {
		return fmt.Errorf("no links found")
	}

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

func init() {
	rootCmd.AddCommand(fileCmd)
}
