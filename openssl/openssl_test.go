package openssl

import (
	"encoding/base64"
	"fmt"
	"os"
	"testing"
)

// TestRun ...
func TestRun(t *testing.T) {
	b, err := Run("rand", "16")
	str := base64.RawURLEncoding.EncodeToString(b)
	t.Log(str, err)
}

// TestEncodeToFile ...
func TestEncodeToFile(t *testing.T) {
	err := KeyToFile("keyfile.key")
	t.Log(err)
}

// TestDecodeFromFile ...
func TestDecodeFromFile(t *testing.T) {
	b, err := DecodeFromFile("keyfile.key")

	t.Log(base64.RawURLEncoding.EncodeToString(b))
	t.Log(err)

	file, err := os.OpenFile("file.key", os.O_RDWR, os.ModePerm)
	p := make([]byte, 1024)

	n, err := file.Read(p)
	p = p[:n]

	s := fmt.Sprintf("%02x", p)
	t.Log(s)

}

// TestFileCount ...
func TestFileCount(t *testing.T) {
	t.Log(FileCount(".", "*.*"))
}
