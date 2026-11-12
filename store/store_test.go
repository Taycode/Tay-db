package store

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestValuesReadBackBeforeAndAfterReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data.log")

	s := open(t, path)
	set(t, s, "user:1", "Tay")
	set(t, s, "user:2", "commas, and\nnewlines")
	set(t, s, "user:1", "Taycode")
	assertValues(t, s, map[string]string{"user:1": "Taycode", "user:2": "commas, and\nnewlines"})
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}

	s = open(t, path)
	defer s.Close()
	assertValues(t, s, map[string]string{"user:1": "Taycode", "user:2": "commas, and\nnewlines"})

	set(t, s, "user:3", "Ada")
	assertValues(t, s, map[string]string{"user:1": "Taycode", "user:3": "Ada"})

	if _, err := s.Get("user:4"); !errors.Is(err, ErrNotFound) {
		t.Errorf("Get(missing key) error = %v, want ErrNotFound", err)
	}
}

func open(t *testing.T, path string) *Store {
	t.Helper()
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func set(t *testing.T, s *Store, key, value string) {
	t.Helper()
	if err := s.Set(key, []byte(value)); err != nil {
		t.Fatal(err)
	}
}

func assertValues(t *testing.T, s *Store, want map[string]string) {
	t.Helper()
	for key, wantValue := range want {
		got, err := s.Get(key)
		if err != nil {
			t.Fatalf("Get(%q): %v", key, err)
		}
		if string(got) != wantValue {
			t.Errorf("Get(%q) = %q, want %q", key, got, wantValue)
		}
	}
}
