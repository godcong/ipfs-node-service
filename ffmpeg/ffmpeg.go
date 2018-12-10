package openssl

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

// Media ...
type Media struct {
	OutPath         string
	MessageCallback func(map[string]interface{}) error
}

// MediaType ...
type MediaType string

// OutPath ...
const (
	OutPath MediaType = "outpath"
)

// NewFFMpeg ...
func NewFFMpeg(args map[MediaType]string) *Media {
	return &Media{
		OutPath: args[OutPath],
	}
}

// Run ...
func (m *Media) Run() {

}

// Run ...
func Run(args ...string) ([]byte, error) {
	if args == nil {
		args = []string{"-h"}
	}
	cmd := exec.Command("ffmpeg", args...)
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

// VideoToMp4 ...
func VideoToMp4(path string, out string) ([]byte, error) {
	return Run("-i", path, "-acodec", "copy", "-vcodec", "copy", out)
}
