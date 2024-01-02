# BIG (batch image getter)

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

## Why BIG is written in Go?

- Go is a compiled language, so it's easy to distribute the binary and run it on any machine
- As the number of images to download can be large, Go is a good choice to handle concurrency
- Go is statically typed so it's easier to write a robust program which is good fit for a command line tool as there no easy way to update it once it's distributed
