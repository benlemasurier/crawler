package main

import "net/url"

// Page represents a page within a website.
type Page struct {
	// relative path to site root e.g., /users/foo
	Path string

	// static assets required by the page
	assets map[string]*url.URL

	// links within the page
	links map[string]*url.URL
}

// Assets returns all assets within a Page.
func (p *Page) Assets() map[string]*url.URL {
	return p.assets
}

// AddAsset inserts a new asset into a Page.
// Duplicate and invalid URLs will be ignored.
func (p *Page) AddAsset(link string) {
	parsed, err := url.Parse(link)
	if err != nil {
		return
	}

	p.assets[parsed.String()] = parsed
}

// Links returns all links within a page.
func (p *Page) Links() map[string]*url.URL {
	return p.links
}

// AddLink inserts a new link into a Page.
// Duplicate, invalid, and blacklisted URLs will be ignored.
func (p *Page) AddLink(link string) {
	parsed, err := url.Parse(link)
	if err != nil {
		return
	}

	// ignore blacklisted schemes
	switch parsed.Scheme {
	case "javascript":
		return
	case "mailto":
		return
	}

	p.links[parsed.String()] = parsed
}

// Normalize modifies a Page's Links ensuring all links with relative
// paths are converted to full and proper urls based on the provided reference.
func (p *Page) Normalize(ref *url.URL) {
	for k, link := range p.Links() {
		if !link.IsAbs() {
			p.Links()[k] = ref.ResolveReference(link)
		}
	}
}

// NewPage returns a new Page with 0 assets and 0 children.
func NewPage(path string) *Page {
	return &Page{
		Path:   path,
		assets: make(map[string]*url.URL),
		links:  make(map[string]*url.URL),
	}
}
