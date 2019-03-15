package openssl

import (
	"encoding/base64"
	"fmt"

	"os"
	"strconv"
	"sync"
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

// TestNumber32 ...
func TestNumber32(t *testing.T) {
	t.Log(Number32(255))
}

// TestRun2 ...
func TestRun2(t *testing.T) {
	file, _ := os.OpenFile("file.key", os.O_RDWR, os.ModePerm)
	p := make([]byte, 1024)

	n, _ := file.Read(p)
	p = p[:n]

	s := fmt.Sprintf("%02x", p)

	wg := sync.WaitGroup{}

	for i := 0; i < 100; i++ {

		wg.Add(1)
		go func(i int) {
			tsFile := "segment-" + strconv.Itoa(i) + ".ts"
			run, e := Run("aes-128-cbc", "-e", "-in", "openssl_test.go", "-out", "encrypted_"+tsFile, "-nosalt", "-iv", Number32(i), "-K", s)
			t.Log(run, e)
			wg.Done()
		}(i)

	}
	wg.Wait()

}

// TestKeyFile ...
func TestKeyFile(t *testing.T) {
	err := KeyFile("tmp", "text", "", "KeyInfoFile", "http://localhost:8080/stream/", true)
	t.Log(err)
}
