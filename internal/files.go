package internal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FileSystem interface {
	CreateDirectory(directory string) error
	Save(body io.ReadCloser, path string) error
	ReadLines(path string) ([]string, error)
	Exists(path string) bool
}

type DefaultFileSystem struct{}

func NewFileSystem() FileSystem {
	return &DefaultFileSystem{}
}

func (f *DefaultFileSystem) CreateDirectory(directory string) error {
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

func (f *DefaultFileSystem) Save(body io.ReadCloser, path string) error {
	defer body.Close()

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, body)
	if err != nil {
		return err
	}

	return nil
}

func (f *DefaultFileSystem) ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)

	var lines []string
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}

	return lines, nil
}

func (f *DefaultFileSystem) Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

const DefaultPath = "."

func validateContentType(contentType string, imageTypes []string) error {
	if !strings.HasPrefix(contentType, "image") {
		return fmt.Errorf("content type '%s' is not an image", contentType)
	}

	imageType := strings.Split(contentType, "/")[1]

	for _, allowedImageType := range imageTypes {
		if imageType == allowedImageType {
			return nil
		}
	}

	return fmt.Errorf("image type '%s' is not allowed", imageType)
}

func addExtensionIfMissing(filePath, contentType string) string {
	extension := filepath.Ext(filePath)

	if extension != "" {
		return filePath
	}

	extension = "." + strings.Split(contentType, "/")[1]

	return filePath + extension
}
