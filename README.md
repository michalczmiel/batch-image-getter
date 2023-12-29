# BIG (batch image downloader)

BIG is a simple commands line tools to download images from a website.

## Usage

For now there is no binary available, so you need to build it yourself.

```bash
go run main.go <url>
```

Specify the extension of the images you want to download with the `-t` flag.

```bash
go run main.go --types png --types jpg <url>
```

Example with real website:

```bash
go run main.go https://www.nasa.gov/image-of-the-day/
```

## Known limitations

- BIG is not able to download images from a website that requires authentication
- BIG is not able to download images from a website that requires JavaScript to display images

## TODO

- [ ] Add a flag to specify the output directory
- [ ] Add a flag to specify the number of concurrent downloads
- [ ] Accept a list of URLs
- [ ] Improve error handling (e.g. when the URL is not valid, when the URL is not reachable, when file is too big etc.)
- [ ] Improve logging (e.g. add a verbose flag)
- [ ] Produce a binary
- [ ] Improve testing coverage
