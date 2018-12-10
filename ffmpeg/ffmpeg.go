package ffmpeg

import (
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

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		log.Println(err)
		return nil, err
	}

	//b, err := ioutil.ReadAll(stdout)
	//if err != nil {
	//	log.Fatal(err)
	//}

	if err := cmd.Wait(); err != nil {
		log.Println(err)
		return nil, err
	}

	return stdout, nil
}

// TranToMp4 ...
func TranToMp4(path string, out string) ([]byte, error) {
	//ffmpeg -i input.mkv -acodec libfaac -vcodec libx264 out.mp4
	return Run("-i", path, "-acodec", "libfaac", "-vcodec", "libx264", out)
}

// CopyToMp4 ...
func CopyToMp4(path string, out string) ([]byte, error) {
	//ffmpeg -i input.mkv -acodec libfaac -vcodec libx264 out.mp4
	return Run("-i", path, "-acodec", "copy", "-vcodec", "copy", out)
}
