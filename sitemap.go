package main

import (
	"fmt"
	"path"
	"sort"
	"strings"
	"sync"
)

// Sitemap represents a site's structure
type Sitemap struct {
	pages map[string]*Page
	sync.Mutex
}

// cleanPath returns the shortest path name equivalent to path.
func cleanPath(p string) string {
	clean := path.Clean(p)
	if clean == "." {
		return "/"
	}

	return clean
}

// Insert adds a new Page to a Sitemap.
// If the Page already exists within the Sitemap it will be updated.
func (s *Sitemap) Insert(p *Page) {
	s.Lock()
	s.pages[cleanPath(p.Path)] = p
	s.Unlock()
}

// Exists determines if path p exists within a sitemap.
func (s *Sitemap) Exists(p string) bool {
	s.Lock()
	_, exists := s.pages[cleanPath(p)]
	s.Unlock()

	return exists
}

// SortedPaths returns all paths within the site sorted alphabetically.
func (s *Sitemap) SortedPaths() []string {
	paths := make([]string, 0, len(s.pages))

	s.Lock()
	for path := range s.pages {
		paths = append(paths, path)
	}
	s.Unlock()

	sort.Strings(paths)
	return paths
}

// Children returns all children of path p.
func (s *Sitemap) Children(p string) []string {
	children := []string{}
	for _, child := range s.SortedPaths() {
		if child != p && strings.HasPrefix(child, p) {
			children = append(children, child)
		}
	}

	return children
}

func (s *Sitemap) String() string {
	var tree string
	printed := make(map[string]bool)

	var printTree func(parent string, chidren []string)
	printTree = func(parent string, children []string) {
		for _, child := range children {
			// skip any previously printed paths
			if _, skip := printed[child]; skip {
				continue
			}

			// show the minimum amount of path leading to basename
			split := strings.Split(child, parent)
			basename := split[len(split)-1]

			// relative path
			indent := len(parent)
			format := fmt.Sprintf("%%%ds/\n", indent+len(basename))
			tree += fmt.Sprintf(format, strings.TrimPrefix(basename, "/"))

			// assets
			if showAssets {
				for _, a := range s.pages[child].Assets() {
					format = fmt.Sprintf("%%%ds/\n", indent+len(basename)+len(a.String())+3)
					tree += fmt.Sprintf(format, " . "+a.String())
				}
			}

			// links
			if showLinks {
				for _, l := range s.pages[child].Links() {
					format = fmt.Sprintf("%%%ds/\n", indent+len(basename)+len(l.String())+3)
					tree += fmt.Sprintf(format, " > "+l.String())
				}
			}

			printed[child] = true
			printTree(child, s.Children(child))
		}
	}

	sorted := s.SortedPaths()
	tree = sorted[0] + "\n"
	printTree(sorted[0], sorted[1:])

	return tree
}

// NewSitemap returns a new (empty) sitemap.
func NewSitemap() *Sitemap {
	return &Sitemap{
		pages: make(map[string]*Page),
	}
}
