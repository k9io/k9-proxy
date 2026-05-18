package main

import (
	"path/filepath"
	"testing"
)

func TestReadCache_Missing(t *testing.T) {
	_, ok := Read_Cache(filepath.Join(t.TempDir(), "missing.cache"))
	if ok {
		t.Fatal("expected false for nonexistent file")
	}
}

func TestWriteCache_SkipsErrors(t *testing.T) {
	file := filepath.Join(t.TempDir(), "test.cache")
	Write_Cache(file, `{"error":"something bad"}`)
	_, ok := Read_Cache(file)
	if ok {
		t.Fatal("error responses must not be cached")
	}
}

func TestWriteReadCache_RoundTrip(t *testing.T) {
	file := filepath.Join(t.TempDir(), "test.cache")
	Write_Cache(file, `{"data":"hello"}`)
	body, ok := Read_Cache(file)
	if !ok {
		t.Fatal("expected true after write")
	}
	if body != `{"data":"hello"}` {
		t.Fatalf("got %q, want %q", body, `{"data":"hello"}`)
	}
}

func TestWriteCache_UpdatesStale(t *testing.T) {
	file := filepath.Join(t.TempDir(), "test.cache")
	Write_Cache(file, `{"data":"v1"}`)
	Write_Cache(file, `{"data":"v2"}`)
	body, ok := Read_Cache(file)
	if !ok {
		t.Fatal("expected true")
	}
	if body != `{"data":"v2"}` {
		t.Fatalf("got %q, want %q", body, `{"data":"v2"}`)
	}
}
