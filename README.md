# Getty

![Build](https://github.com/m-rots/getty/workflows/Build/badge.svg)

A small utility to download images from gettyimages.
- Images are locked to 2048x2048 pixels in size

## Usage
Download an image by its identifier:
```bash
getty "186350194"
```

Download an image by its URL:
```bash
getty "https://www.gettyimages.com/detail/186350194"
```

Download multiple images:
```bash
getty "186350194" "1207089198"
```

Downloaded images will be placed in the current working directory with the identifier of the image as a JPG file.

## Install
Either download one of the binaries in the releases tab or build getty yourself.

## Build
1. Make sure you have the latest version of [Go](https://golang.org/dl/) installed
2. Run `go build getty.go` to build the binary