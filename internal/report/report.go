package report

import (
	"fmt"
	"sort"
)

type Page struct {
	URL   string
	Count int
}

func Print(pages map[string]int, baseURL string) {
	fmt.Printf(`
======================================
 REPORT for %s
======================================
`, baseURL)

	sortedPages := sortPages(pages)
	for _, page := range sortedPages {
		fmt.Printf("Found %d internal links to %s\n", page.Count, page.URL)
	}
}

func sortPages(pages map[string]int) []Page {
	pairs := []Page{}
	for url, count := range pages {
		pairs = append(pairs, Page{URL: url, Count: count})
	}

	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Count == pairs[j].Count {
			return pairs[i].URL < pairs[j].URL
		}
		return pairs[i].Count > pairs[j].Count
	})

	return pairs
}
