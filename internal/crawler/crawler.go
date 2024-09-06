package crawler

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

type Crawler struct {
	BaseURL url.URL
	Pages   map[string]int
	wg      *sync.WaitGroup
	mu      *sync.Mutex

	concurrencyControl chan struct{}
	maxPages           int
}

func New(rawBaseURL string, maxConcurreny, maxPages int) (*Crawler, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		msg := fmt.Sprintf("Couldn't parse base URL: %s", err)
		return nil, errors.New(msg)
	}

	c := &Crawler{
		BaseURL:            *baseURL,
		Pages:              make(map[string]int),
		wg:                 &sync.WaitGroup{},
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurreny),
		maxPages:           maxPages,
	}
	return c, nil
}

func (c *Crawler) Start() {
	c.wg.Add(1)
	c.crawlPage(c.BaseURL.String())
	c.wg.Wait()
}

func (c *Crawler) crawlPage(rawCurrentURL string) {
	c.concurrencyControl <- struct{}{}
	defer func() {
		<-c.concurrencyControl
		c.wg.Done()
	}()

	if len(c.Pages) >= c.maxPages {
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		log.Printf("Couldn't parse current url: %s\n", err)
		return
	}

	if currentURL.Hostname() != c.BaseURL.Hostname() {
		return
	}

	normalizeURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		log.Printf("Couldn't normalized currentURL: %s\n", err)
	}

	if isFirst := c.addPageVisit(normalizeURL); !isFirst {
		return
	}

	log.Printf("Crawling - %s\n", rawCurrentURL)

	htmlBody, err := getHTML(rawCurrentURL)
	if err != nil {
		log.Printf("Error getHTML: %s\n", err)
		return
	}

	nextURLs, err := getURLsFromHTMLBody(htmlBody, rawCurrentURL)
	if err != nil {
		log.Printf("Error getURLsFromHTML: %s'n", err)
		return
	}

	for _, nextURL := range nextURLs {
		c.wg.Add(1)
		go c.crawlPage(nextURL)
	}
}

func (c *Crawler) addPageVisit(normalizeURL string) (isFirst bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, visited := c.Pages[normalizeURL]; visited {
		c.Pages[normalizeURL]++
		return false
	}

	c.Pages[normalizeURL] = 1
	return true
}

func getURLsFromHTMLBody(htmlBody, rawBaseURL string) ([]string, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse base URL: %v", err)
	}

	doc, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return nil, fmt.Errorf("Couldn't parse HTML: %s", err)
	}

	var urls []string
	var traverseNodes func(node *html.Node)
	traverseNodes = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, a := range node.Attr {
				if a.Key == "href" {
					href, err := url.Parse(a.Val)
					if err != nil {
						log.Printf("Couldn't parse href %s: %s\n", a.Val, err)
						continue
					}

					resolvedURL := baseURL.ResolveReference(href)
					urls = append(urls, resolvedURL.String())
				}
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			traverseNodes(c)
		}
	}

	traverseNodes(doc)
	return urls, nil
}

func getHTML(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", fmt.Errorf("Couldn't make a request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		return "", fmt.Errorf("HTTP error: %s", resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("Got non HTML response %s", contentType)
	}

	htmlBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Couldn't read resp.Body: %s", err)
	}

	return string(htmlBody), nil
}

func normalizeURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("Couldn't normalize URL: %s", err)
	}

	newURL := parsedURL.Host + parsedURL.Path
	newURL = strings.ToLower(newURL)
	newURL = strings.TrimSuffix(newURL, "/")

	return newURL, nil
}
