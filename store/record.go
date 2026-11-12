package store

import (
	"encoding/binary"
	"io"
)

// Every record in the log is laid out as:
//
//	| key size (4 bytes) | value size (4 bytes) | key | value |
//
// The sizes come first so a reader always knows how many bytes to read next.
const headerSize = 8

func encodeRecord(key string, value []byte) []byte {
	record := make([]byte, headerSize+len(key)+len(value))
	binary.LittleEndian.PutUint32(record[0:4], uint32(len(key)))
	binary.LittleEndian.PutUint32(record[4:8], uint32(len(value)))
	copy(record[headerSize:], key)
	copy(record[headerSize+len(key):], value)
	return record
}

// readHeader returns io.EOF only when the log ends cleanly between records.
func readHeader(r io.Reader) (keySize, valueSize uint32, err error) {
	header := make([]byte, headerSize)
	if _, err := io.ReadFull(r, header); err != nil {
		return 0, 0, err
	}
	return binary.LittleEndian.Uint32(header[0:4]), binary.LittleEndian.Uint32(header[4:8]), nil
}
