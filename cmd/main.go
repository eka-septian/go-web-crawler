package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/eka-septian/go-web-crawler/internal/crawler"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("usage: ./crawler <URL string> <maxConcurrency int> <maxPages int>")
		os.Exit(1)
	}

	rawBaseURL := os.Args[1]
	maxConcurrent, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("usage: ./crawler <URL string> <maxConcurrency int> <maxPages int>")
		os.Exit(1)
	}
	maxPages, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("usage: ./crawler <URL string> <maxConcurrency int> <maxPages int>")
		os.Exit(1)
	}

	c, err := crawler.New(rawBaseURL, maxConcurrent, maxPages)
	if err != nil {
        panic(err)
	}
	c.Start()
}
