package openssl

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
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

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}

	return b, nil
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
func KeyFile(path, fname string, key, uri string, iv bool) error {
	var err error

	//newKey := KeyToHex(key)
	err = SaveTo(path+fname+"/key", key)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path+fname+"/KeyInfo", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	_, _ = file.WriteString(uri + "/" + fname + "/key")
	_, _ = file.WriteString("\n")
	_, _ = file.WriteString(path + fname + "/key")
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
