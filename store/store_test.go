package store

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestReopenRebuildsIndexFromLog(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data.log")

	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	writes := []struct{ key, value string }{
		{"user:1", "Tay"},
		{"user:2", "commas, and\nnewlines"},
		{"user:1", "Taycode"},
	}
	for _, w := range writes {
		if err := s.Set(w.key, []byte(w.value)); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}

	s, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	want := map[string]string{"user:1": "Taycode", "user:2": "commas, and\nnewlines"}
	for key, wantValue := range want {
		got, err := s.Get(key)
		if err != nil {
			t.Fatalf("Get(%q) after reopen: %v", key, err)
		}
		if string(got) != wantValue {
			t.Errorf("Get(%q) after reopen = %q, want %q", key, got, wantValue)
		}
	}
	if _, err := s.Get("user:3"); !errors.Is(err, ErrNotFound) {
		t.Errorf("Get(missing key) error = %v, want ErrNotFound", err)
	}
}
