# Go Web Crawler
CLI application that generates an "internal links" report for any website on the internet by crawling each page of the site.
This project is created as a way for me to learn and practice Golang

## How to Use
To run this application, you'll need [Git](https://git-scm.com) and [Go](https://go.dev/dl/) installed on your computer. From your command line:

```bash
# Clone this repository
$ git clone https://github.com/ekastn/go-web-crawler

# Go into the repository
$ cd go-web-crawler

# Build the application
$ go build -o bin/crawler cmd/main.go

# Run the crawler
$ ./bin/crawler <URL string> <maxConcurrency int> <maxPages int>

# Example
$ ./bin/crawler https://crawler-test.com 10 50
```
