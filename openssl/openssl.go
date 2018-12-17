package openssl

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// OpenSSL ...
type OpenSSL struct {
	key []byte
}

// Run ...
func Run(args ...string) ([]byte, error) {
	if args == nil {
		args = []string{"help"}
	}
	cmd := exec.Command("openssl", args...)
	cmd.Env = os.Environ()
	cmd.Env = os.Environ()

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return stdout, err
	}

	return stdout, nil
}

// KeyToFile ...
func KeyToFile(path string) error {
	key, err := Run("rand", "-base64", "20")
	if err != nil {
		return err
	}

	return EncodeToFile(key, path)
}

// Base64Key ...
func Base64Key() ([]byte, error) {
	return Run("rand", "-base64", "20")
}

// HexIV ...
func HexIV() ([]byte, error) {
	return Run("rand", "-hex", "16")
}

// HexKey ...
func HexKey() ([]byte, error) {
	return Run("rand", "-hex", "32")
}

// KeyToHex ...
func KeyToHex(key []byte) string {
	return fmt.Sprintf("%02x", key)
}

// SaveTo ...
func SaveTo(path string, data string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = file.WriteString(data)
	return err
}

// EncodeToFile ...
func EncodeToFile(key []byte, path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}

	enc := base64.NewEncoder(base64.RawURLEncoding, file)
	_, err = enc.Write(key)

	if err != nil {
		return err
	}
	return nil
}

// DecodeFromFile ...
func DecodeFromFile(path string) ([]byte, error) {
	file, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}
	p := make([]byte, 1024)

	dec := base64.NewDecoder(base64.RawURLEncoding, file)
	i, err := dec.Read(p)

	if err != nil {
		return nil, err
	}
	return p[:i], nil
}

// FileCount ...
func FileCount(path, name string) int {
	infos, err := filepath.Glob(path + "/" + name)
	if err != nil {
		return 0
	}
	return len(infos)
}

// Number32 ...
func Number32(i int) string {
	return fmt.Sprintf("%032x", i)
}

// KeyFile ...
func KeyFile(path, keyName, key, keyInfo, uri string, iv bool) error {
	var err error

	err = SaveTo(path+"/"+keyName, key)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path+"/"+keyInfo, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	_, _ = file.WriteString(uri)
	_, _ = file.WriteString("\n")
	_, _ = file.WriteString(path + "/" + keyName)
	_, _ = file.WriteString("\n")
	if iv {
		key, err := Run("rand", "-hex", "16")
		if err != nil {
			return err
		}
		_, _ = file.Write(key)
	}

	return nil
}
