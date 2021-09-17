package utils

import (
	"bytes"
	"compress/gzip"
)

var (
	gzipIdent = []byte("\x1F\x8B\x08")
	wasmIdent = []byte("\x00\x61\x73\x6D")
)

// IsGzip returns checks if the file contents are gzip compressed
func IsGzip(input []byte) bool {
	return bytes.Equal(input[:3], gzipIdent)
}

// IsWasm checks if the file contents are of wasm binary
func IsWasm(input []byte) bool {
	return bytes.Equal(input[:4], wasmIdent)
}

// GzipIt compresses the input ([]byte)
func GzipIt(input []byte) ([]byte, error) {
	// Create gzip writer.
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(input)
	if err != nil {
		return nil, err
	}
	err = w.Close() // You must close this first to flush the bytes to the buffer.
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
