package store

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

var ErrNotFound = errors.New("key not found")

// Store is an append-only log on disk plus an in-memory index of where each key's newest value sits.
type Store struct {
	file  *os.File
	index map[string]location
	end   int64
}

type location struct {
	offset int64
	size   uint32
}

func Open(path string) (*Store, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open log %s: %w", path, err)
	}

	s := &Store{file: file, index: make(map[string]location)}
	if err := s.rebuildIndex(); err != nil {
		file.Close()
		return nil, fmt.Errorf("rebuild index from %s: %w", path, err)
	}
	return s, nil
}

func (s *Store) Set(key string, value []byte) error {
	record := encodeRecord(key, value)
	if _, err := s.file.Write(record); err != nil {
		return fmt.Errorf("append %q: %w", key, err)
	}

	s.index[key] = location{offset: s.end + headerSize + int64(len(key)), size: uint32(len(value))}
	s.end += int64(len(record))
	return nil
}

func (s *Store) Get(key string) ([]byte, error) {
	loc, ok := s.index[key]
	if !ok {
		return nil, ErrNotFound
	}

	value := make([]byte, loc.size)
	if _, err := s.file.ReadAt(value, loc.offset); err != nil {
		return nil, fmt.Errorf("read %q at offset %d: %w", key, loc.offset, err)
	}
	return value, nil
}

func (s *Store) Len() int {
	return len(s.index)
}

func (s *Store) Close() error {
	return s.file.Close()
}

// rebuildIndex reads the whole log from the start. A later record for the same key
// overwrites the earlier index entry, so each key ends up pointing at its newest value.
func (s *Store) rebuildIndex() error {
	r := bufio.NewReader(s.file)
	var offset int64

	for {
		keySize, valueSize, err := readHeader(r)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read record header at offset %d: %w", offset, err)
		}

		key := make([]byte, keySize)
		if _, err := io.ReadFull(r, key); err != nil {
			return fmt.Errorf("read key at offset %d: %w", offset, err)
		}
		if _, err := r.Discard(int(valueSize)); err != nil {
			return fmt.Errorf("skip value of %q at offset %d: %w", key, offset, err)
		}

		valueOffset := offset + headerSize + int64(keySize)
		s.index[string(key)] = location{offset: valueOffset, size: valueSize}
		offset = valueOffset + int64(valueSize)
	}

	s.end = offset
	return nil
}
