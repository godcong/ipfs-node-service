package openssl

import (
	"encoding/base64"
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
}
