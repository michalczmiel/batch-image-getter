# BIG (batch image getter)

BIG is a simple commands line tools for downloading images.

## Usage

For now there is no binary available, so you need to build it yourself.

### Fetching images from a HTML page

```bash
go run main.go html <url>
```

Example with real website:

```bash
go run main.go html https://www.nasa.gov/image-of-the-day/ -d space-images
```

### Fetching images from text file containing URLs

```bash
go run main.go file <path>
```

### Customization using flags

```bash
Flags:
  -c, --concurrency int          number of concurrent downloads (default 10)
  -d, --dir string               directory to save images to (default ".")
  -h, --help                     help for file
      --json                     output results as json
      --max-sleep-interval int   maximum number of seconds to sleep after each request
      --referer string           custom referer to use for requests
      --sleep-interval int       number of seconds to sleep after each request or minimum number of seconds for randomized sleep when used along max-sleep-interval (default 0)
      --types stringArray        image types to download (default [jpg,jpeg,png,gif,webp])
      --user-agent string        custom user agent to use for requests
```

### Running tests

For Go unit tests:

```bash
go test ./...
```

for Python integration tests:

```bash
pytest
```

## Known limitations

- BIG is not able to download images from a website that requires authentication
- BIG is not able to download images from a website that requires JavaScript to display images

## Why BIG is written in Go?

- Go is a compiled language, so it's easy to distribute the binary and run it on any machine
- As the number of images to download can be large, Go is a good choice to handle concurrency
- Go is statically typed so it's easier to write a robust program which is good fit for a command line tool as there no easy way to update it once it's distributed
