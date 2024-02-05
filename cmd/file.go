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

	fileSystem := internal.NewFileSystem()

	if !fileSystem.Exists(path) {
		return fmt.Errorf("file does not exist")
	}

	return nil
}

var fileCmd = &cobra.Command{
	Use:   "file <path>",
	Short: "Download all images from a file",
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

	fileSystem := internal.NewFileSystem()
	httpClient := internal.NewHttpClient(parameters.UserAgent)

	provider := internal.NewFileProvider(filePath, fileSystem)

	links, err := provider.Links()
	if err != nil {
		return err
	}

	return internal.Run(links, parameters, httpClient, fileSystem)
}

func init() {
	rootCmd.AddCommand(fileCmd)
}
