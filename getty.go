package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

const usage string = `Getty ya pirate!

Usage: getty <id | url>
`

var client = &http.Client{}

func downloadJpg(url string, img chan image.Image) {
	res, err := client.Get(url)
	if err != nil || res.StatusCode != 200 {
		fmt.Println("Could not download the image")
		os.Exit(1)
	}

	defer res.Body.Close()

	data, err := jpeg.Decode(res.Body)
	if err != nil {
		fmt.Println("Could not decode the image")
		os.Exit(1)
	}

	img <- data
}

func mergeImages(imageOne image.Image, imageTwo image.Image) image.Image {
	x, y := imageOne.Bounds().Dx(), imageOne.Bounds().Dy()

	rectOne := image.Rect(0, 0, x, y/5*3)
	rectTwo := image.Rect(0, y/5*3, x, y)

	newImage := image.NewRGBA(imageOne.Bounds())

	draw.Draw(newImage, rectOne, imageOne, image.ZP, draw.Src)
	draw.Draw(newImage, rectTwo, imageTwo, image.Pt(0, y/5*3), draw.Src)

	return newImage
}

func getty(id string, wg *sync.WaitGroup) {
	defer wg.Done()

	urlOne := fmt.Sprintf("https://media.gettyimages.com/photos/-id%s?s=2048x2048&w=5", id)
	urlTwo := fmt.Sprintf("https://media.gettyimages.com/photos/-id%s?s=2048x2048&w=125", id)

	channelOne := make(chan image.Image)
	channelTwo := make(chan image.Image)
	go downloadJpg(urlOne, channelOne)
	go downloadJpg(urlTwo, channelTwo)

	newImage := mergeImages(<-channelOne, <-channelTwo)

	resultPath := fmt.Sprintf("%s.jpg", id)
	resultFile, err := os.Create(resultPath)
	if err != nil {
		fmt.Println("Could not write the resulting image")
		os.Exit(1)
	}

	jpeg.Encode(resultFile, newImage, &jpeg.Options{
		Quality: 100,
	})
	defer resultFile.Close()
}

func main() {
	args := os.Args[1:]

	help := false

	if len(args) == 0 || args[0] == "help" {
		help = true
	}

	var links []string

	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			help = true
		} else if arg[0] == '-' {
			fmt.Println("Unknown argument " + arg)
			os.Exit(1)
		} else {
			links = append(links, arg)
		}
	}

	if help {
		fmt.Println(usage)
		return
	}

	var wg sync.WaitGroup

	for _, link := range links {
		var id string
		if isGettyID(link) {
			id = link
		} else {
			id = idFromURL(link)
		}

		wg.Add(1)
		go getty(id, &wg)
	}

	wg.Wait()
}

func idFromURL(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		fmt.Println("Invalid url: " + link)
		os.Exit(1)
	}

	parts := strings.Split(u.Path, "/")
	id := parts[len(parts)-1]

	if !isGettyID(id) {
		fmt.Println("Invalid url: " + link)
		os.Exit(1)
	}

	return id
}

func isGettyID(id string) bool {
	if _, err := strconv.Atoi(id); err == nil {
		return true
	}

	return false
}
