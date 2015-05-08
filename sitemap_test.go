package main

import "testing"

func TestInsert(t *testing.T) {
	s := NewSitemap()
	p := NewPage("/test")

	s.Insert(p)
	if !s.Exists(p.Path) {
		t.Errorf("expected page to be inserted")
	}
}

func TestExists(t *testing.T) {
	s := NewSitemap()
	p := NewPage("/test")
	s.Insert(p)

	if !s.Exists(p.Path) {
		t.Errorf("expected exists to return true")
	}
}

func TestSortedPaths(t *testing.T) {
	s := NewSitemap()
	s.Insert(NewPage("/bar"))
	s.Insert(NewPage("/foo"))
	s.Insert(NewPage("/"))

	expected := []string{"/", "/bar", "/foo"}
	sorted := s.SortedPaths()
	for i := 0; i < len(sorted); i++ {
		if sorted[i] != expected[i] {
			t.Errorf("expected %s to equal %s", sorted[i], expected[i])
		}
	}
}

func TestChildren(t *testing.T) {
	s := NewSitemap()
	s.Insert(NewPage("/foo"))
	s.Insert(NewPage("/"))

	expected := "/foo"
	actual := s.Children("/")[0]
	if expected != actual {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestString(t *testing.T) {
	showAssets = true
	showLinks = true

	s := NewSitemap()
	foo := NewPage("/foo")
	foo.AddLink("https://example.com")
	bar := NewPage("/foo/bar")
	bar.AddAsset("//cdn.example.com/hipster.js")

	s.Insert(NewPage("/"))
	s.Insert(foo)
	s.Insert(bar)

	expected := "/\n" +
		" foo/\n" +
		"     > https://example.com/\n" +
		"     bar/\n" +
		"         . //cdn.example.com/hipster.js/\n"

	actual := s.String()
	if actual != expected {
		t.Errorf("expected \n%s\ngot\n%s\n", expected, actual)
	}
}
