package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// Fetcher is a Page fetching interface.
type Fetcher interface {
	Fetch(link *url.URL) (*Page, error)
}

// HTTPFetcher fetches Pages over HTTP(S).
type HTTPFetcher struct{}

// Fetch performs an HTTP GET on the provided URL and
// returns a Page populated by parsing the returned HTML.
func (f HTTPFetcher) Fetch(link *url.URL) (*Page, error) {
	fmt.Println("fetching:", link)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: skipVerify},
		},
	}
	resp, err := client.Get(link.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("received http %d for %s", resp.StatusCode, link)
	}

	p := NewPage(resp.Request.URL.Path)
	if err := parse(resp.Body, p); err != nil {
		return nil, err
	}

	// resolve relative urls
	p.Normalize(link)

	return p, err
}

// parse traverses a html document populating Page p's assets and links.
func parse(r io.Reader, p *Page) error {
	doc, err := html.Parse(r)
	if err != nil {
		return err
	}

	var search func(n *html.Node)
	search = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "a":
				// ignore empty and anchor links
				link := extract(n, "href")
				if strings.TrimSpace(strings.Split(link, "#")[0]) != "" {
					p.AddLink(link)
				}
			case "link":
				p.AddAsset(extract(n, "href"))
			case "script":
				p.AddAsset(extract(n, "src"))
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			search(c)
		}
	}

	search(doc)

	return nil
}

// pulls the value for key from an html node.
func extract(n *html.Node, key string) string {
	for _, tag := range n.Attr {
		if tag.Key == key {
			return tag.Val
		}
	}

	return ""
}
