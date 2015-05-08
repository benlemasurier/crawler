package main

import (
	"fmt"
	"net/url"
	"testing"
)

var testPages = map[string]bool{
	"https://example.com":                                true,
	"https://example.com/":                               true,
	"https://assets.example.com":                         true,
	"https://example.com/foo":                            true,
	"https://example.com/foo/bar/baz":                    false,
	"https://example.com/foo/bar/baz/whahappenafterbaz?": false,
	"/bar":     true,
	"/foo/bar": true,
}

type testFetcher map[string]*Page

func (tf testFetcher) Fetch(link *url.URL) (*Page, error) {
	if valid, exists := testPages[link.String()]; valid && exists {
		p := NewPage(link.Path)
		p.AddAsset("//cdn.example.com/jquery.js")
		for k := range testPages {
			p.AddLink(k)
		}

		return p, nil
	}

	return nil, fmt.Errorf("%s not found", link)
}

var tf = make(testFetcher)

func TestCrawl(t *testing.T) {
	u, _ := url.Parse("https://example.com")
	siteURL = u

	s := NewSitemap()
	crawl(siteURL, tf, s)
}

func TestExternalLink(t *testing.T) {
	u, _ := url.Parse("https://example.com")
	siteURL = u

	test, _ := url.Parse("http://assets.example.com")
	if !external(test) {
		t.Logf("expected %s to be an external link", test)
	}
}
