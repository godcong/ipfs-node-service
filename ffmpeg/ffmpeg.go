package ffmpeg

import (
	"log"
	"os"
	"os/exec"
)

// Mpeg ...
type Mpeg struct {
	OutPath         string
	MessageCallback func(map[string]interface{}) error
}

// MpegType ...
type MpegType string

// OutPath ...
const (
	OutPath MpegType = "outpath"
)

// New ...
func New(args map[MpegType]string) *Mpeg {
	return &Mpeg{
		OutPath: args[OutPath],
	}
}

// Run ...
func (m *Mpeg) Run() {

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
		//log.Println(string(stdout), err)
		return stdout, err
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
