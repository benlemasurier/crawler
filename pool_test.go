package main

import "testing"

func TestNewPool(t *testing.T) {
	workers := 5
	p := NewPool(workers)
	if cap(p) != workers {
		t.Errorf("expected a pool of %d workers\n", workers)
	}
}
