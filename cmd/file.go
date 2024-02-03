package cmd

import (
	"fmt"

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

	if !internal.DoesFileExist(path) {
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

		err = validateRootFlags(cmd)
		if err != nil {
			return err
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

	lines, err := internal.GetLinesFromFile(filePath)
	if err != nil {
		return err
	}

	var links []string
	for _, line := range lines {
		if internal.IsUrlValid(line) {
			links = append(links, line)
		}
	}

	links = internal.RemoveDuplicates(links)

	if len(links) == 0 {
		return fmt.Errorf("no links found")
	}

	err = internal.CreateDirectoryIfDoesNotExists(parameters.Directory)
	if err != nil {
		return err
	}

	httpClient := internal.NewHttpClient()

	inputs := internal.PrepareLinksForDownload(links, parameters)
	results := internal.DownloadImages(inputs, httpClient, parameters)

	printer := internal.NewStdoutPrinter()
	printer.PrintResultsAsPlainText(results)

	return nil
}

func init() {
	rootCmd.AddCommand(fileCmd)
}
