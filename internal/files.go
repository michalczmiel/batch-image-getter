package internal

import (
	"fmt"
	"os"
)

const DefaultPath = "."

func CreateDirectoryIfDoesNotExists(directory string) error {
	if directory == DefaultPath {
		return nil
	}

	if directory == "" {
		directory = DefaultPath
	}

	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directory %s %v", directory, err)
	}

	return nil
}
