package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"testing/iotest"
)

func setup(t *testing.T) *os.File {
	f, err := os.Open("test/test.html")
	if err != nil {
		t.Fatal(err)
	}

	return f
}

func server(f io.Reader) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, f)
	}))
}

func notFoundServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
}

func TestFetch(t *testing.T) {
	f := setup(t)
	defer f.Close()

	s := server(f)
	defer s.Close()

	link, err := url.Parse(s.URL)
	if err != nil {
		t.Fatal(err)
	}

	var fetcher HTTPFetcher
	p, err := fetcher.Fetch(link)
	if err != nil {
		t.Error(err)
	}

	if len(p.Links()) < 1 {
		t.Error("expected page to have links")
	}

	if len(p.Assets()) < 1 {
		t.Error("expected page to have assets")
	}
}

func TestFetch404(t *testing.T) {
	s := notFoundServer()
	defer s.Close()

	link, err := url.Parse(s.URL)
	if err != nil {
		t.Fatal(err)
	}

	var fetcher HTTPFetcher
	_, err = fetcher.Fetch(link)
	if err == nil {
		t.Errorf("expected fetch to fail on 404")
	}
}

func TestParseInvalid(t *testing.T) {
	buf := iotest.TimeoutReader(bytes.NewBufferString("alsdfkjasf"))
	if err := parse(buf, nil); err == nil {
		t.Errorf("expected invalid html to cause a parse failure")
	}
}
