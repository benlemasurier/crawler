package main

import (
	"flag"
	"net/url"
)

var (
	// site root url (flag)
	site = "https://example.com"

	// crawl root URL
	siteURL *url.URL

	// number of concurrent requests
	workers = 10

	// show assets in sitemap output?
	showAssets = false

	// show links in sitemap output?
	showLinks = false

	// disabled tls verification
	skipVerify = false
)

func init() {
	flag.BoolVar(&showAssets, "assets", showAssets, "show page assets in sitemap output")
	flag.BoolVar(&showLinks, "links", showLinks, "show page links in sitemap output")
	flag.IntVar(&workers, "concurrency", workers, "number of concurrent requests")
	flag.StringVar(&site, "url", site, "url to crawl")
	flag.BoolVar(&skipVerify, "insecure", skipVerify, "ignore invalid site certificates")
}
