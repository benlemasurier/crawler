package main

import (
	"fmt"
	"net/url"
	"testing"
)

func TestNewPage(t *testing.T) {
	path := "/test"
	p := NewPage(path)

	// path should match what we passed in NewPath()
	if p.Path != path {
		t.Error("expected page path to match")
	}

	// assets should be empty
	if len(p.Assets()) > 0 {
		t.Error("expected new page to have 0 assets")
	}

	// links should be empty
	if len(p.Links()) > 0 {
		t.Error("expected new page to have 0 links")
	}
}

func TestAssets(t *testing.T) {
	p := NewPage("/test")

	p.AddAsset("http://example.com")

	if len(p.Assets()) != 1 {
		t.Errorf("expected 1 asset to be returned")
	}
}

func TestAddAsset(t *testing.T) {
	p := NewPage("/test")

	p.AddAsset("http://example.com")
	if _, found := p.Assets()["http://example.com"]; !found {
		t.Error("expected asset to be added")
	}
}

func TestAddAssetInvalid(t *testing.T) {
	p := NewPage("/test")

	p.AddAsset(":")
	if _, found := p.Assets()[":"]; found {
		t.Error("expected asset to be invalid")
	}
}

func TestLinks(t *testing.T) {
	p := NewPage("/test")

	p.AddLink("http://example.com")

	if len(p.Links()) != 1 {
		t.Errorf("expected 1 asset to be returned")
	}
}

func TestAddLink(t *testing.T) {
	p := NewPage("/test")

	p.AddLink("http://example.com")
	if _, found := p.Links()["http://example.com"]; !found {
		t.Error("expected asset to be added")
	}
}

func TestAddLinkBlacklisted(t *testing.T) {
	p := NewPage("/test")

	blacklisted := []string{
		"javascript:void(0)",
		"mailto:test@example.com",
	}

	for _, link := range blacklisted {
		p.AddLink(link)
	}

	for link := range p.Links() {
		for _, fail := range blacklisted {
			if link == fail {
				t.Errorf("expected page not to add blacklisted scheme: '%s'", fail)
			}
		}
	}
}

func TestAddLinkInvalid(t *testing.T) {
	p := NewPage("/test")

	p.AddLink(":")
	if _, found := p.Links()[":"]; found {
		t.Error("expected link to be invalid")
	}
}

func TestNormalize(t *testing.T) {
	p := NewPage("/")

	p.AddLink("/test")

	ref, _ := url.Parse("http://example.com")
	p.Normalize(ref)

	if p.Links()["/test"].String() != "http://example.com/test" {
		fmt.Println(p.Links()["/test"])
		t.Error("expected Normalize to add link host and scheme")
	}
}
