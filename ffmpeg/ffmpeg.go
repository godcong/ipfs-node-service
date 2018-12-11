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
func Run(args ...string) (string, error) {
	if args == nil {
		args = []string{"-h"}
	}
	cmd := exec.Command("ffmpeg", args...)
	log.Println(cmd.Args)
	cmd.Env = os.Environ()

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		//log.Println(string(stdout), err)
		return string(stdout), err
	}

	//if err := cmd.Start(); err != nil {
	//	log.Println(err)
	//	return nil, err
	//}

	//b, err := ioutil.ReadAll()
	//if err != nil {
	//	log.Println(string(b))
	//}

	//if err := cmd.Wait(); err != nil {
	//	log.Println(err)
	//	return nil, err
	//}

	return string(stdout), nil
}

// TranToMp4 ...
func TranToMp4(path string, out string) (string, error) {
	//ffmpeg -i input.mkv -acodec libfaac -vcodec libx264 out.mp4
	return Run("-i", path, "-y", "-vcodec", "libx264", "-acodec", "aac", out)

}

// CopyToMp4 ...
func CopyToMp4(path string, out string) (string, error) {
	//cmd:ffmpeg -i input.mkv -acodec copy -vcodec copy out.mp4
	return Run("-i", path, "-y", "-acodec", "copy", "-vcodec", "copy", out)
}

// TransToTS ...
func TransToTS(path string, out string) (string, error) {
	//cmd:ffmpeg -i INPUT.mp4 -codec copy -bsf:v h264_mp4toannexb OUTPUT.ts
	return Run("-i", path, "-y", "-codec", "copy", "-bsf:v", "h264_mp4toannexb", out)
}
