package provider

import "github.com/michalczmiel/batch-image-getter/internal"

type FileProvider struct {
	path       string
	fileSystem internal.FileSystem
}

func NewFileProvider(path string, fileSystem internal.FileSystem) Provider {
	return &FileProvider{
		path:       path,
		fileSystem: fileSystem,
	}
}

func (p *FileProvider) Links() ([]string, error) {
	lines, err := p.fileSystem.ReadLines(p.path)
	if err != nil {
		return nil, err
	}

	var links []string
	for _, line := range lines {
		if internal.IsUrlValid(line) {
			links = append(links, line)
		}
	}

	return links, nil
}
