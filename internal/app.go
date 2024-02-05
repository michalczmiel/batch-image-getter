package internal

func Run(links []string, parameters *Parameters, httpClient HttpClient, fileSystem FileSystem) error {
	links = RemoveDuplicates(links)
	printer := NewStdoutPrinter(parameters.OutputFormat)
	printer.PrintProgress(len(links))

	err := fileSystem.CreateDirectory(parameters.Directory)
	if err != nil {
		return err
	}

	inputs := PrepareLinksForDownload(links, parameters)
	results := DownloadImages(inputs, httpClient, fileSystem, parameters)
	printer.PrintResults(results)

	return nil
}
