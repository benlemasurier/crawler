package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
)

// throttles outgoing connections
var pool = NewPool(workers)

// tracks fetched pages
var fetched = struct {
	urls map[string]bool
	sync.Mutex
}{urls: make(map[string]bool)}

// external determines whether pageURL is external to the global siteURL.
// This is a strict check, domain CNAMES will return false
//	 e.g., example.com != www.example.com
func external(link *url.URL) bool {
	// if the hosts match or this is a relative URL we're good.
	if link.Host == siteURL.Host || link.Host == "" {
		return false
	}

	return true
}

// visited checks whether a given url has been or is currently being fetched.
func visited(u *url.URL) bool {
	fetched.Lock()
	_, seen := fetched.urls[strings.TrimRight(u.String(), "/")]
	fetched.Unlock()

	return seen
}

// mark sets a url as having been fetched.
func mark(u *url.URL) {
	fetched.Lock()
	fetched.urls[strings.TrimRight(u.String(), "/")] = true
	fetched.Unlock()
}

// fetchable returns an array of fetchable URLs.
// URLs are considered fetchable if they are within the crawler's domain.
func fetchable(links map[string]*url.URL) []*url.URL {
	var internal []*url.URL
	for _, link := range links {
		// skip external links
		if external(link) {
			continue
		}

		internal = append(internal, link)
	}

	return internal
}

func crawl(pageURL *url.URL, f Fetcher, s *Sitemap) {
	// skip links we've visited
	if visited(pageURL) {
		return
	}

	// prevent this page from being fetched again
	mark(pageURL)

	// wait for an available worker, fetch, and release worker
	<-pool
	p, err := f.Fetch(pageURL)
	pool <- true
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch error: %s\n", err)
		return
	}

	// add page to sitemap
	s.Insert(p)

	// crawl every local link
	done := make(chan bool)
	links := fetchable(p.Links())
	for _, link := range links {
		go func(u *url.URL) {
			crawl(u, f, s)
			done <- true
		}(link)
	}

	// wait for all requests to finish
	for range links {
		<-done
	}
}

func main() {
	flag.Parse()

	// validate root URL
	parsed, err := url.Parse(site)
	if err != nil {
		fmt.Fprintf(os.Stderr, "crawl url error: %s\n", err)
		os.Exit(1)
	}

	siteURL = parsed
	if siteURL.Scheme == "" {
		siteURL.Scheme = "http"
	}

	s := NewSitemap()
	f := HTTPFetcher{}

	crawl(siteURL, f, s)

	// show tree
	fmt.Println(s)
}
