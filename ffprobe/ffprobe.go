package ffprobe

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

// Probe ...
type Probe struct {
	cmd    []string
	output string
	err    error
}

// New ...
func New(args ...string) *Probe {
	return &Probe{
		cmd:    args,
		output: "",
		err:    nil,
	}
}

// Run ...
func (p *Probe) Run() {
	p.output, p.err = Run(p.cmd...)
	return
}

// Err ...
func (p *Probe) Err() error {
	return p.err
}

// IsH264AndAAC ...
func (p *Probe) IsH264AndAAC() bool {
	video := filterStream(p.output, "Video")
	audio := filterStream(p.output, "Audio")
	if CheckH264(video) && CheckAAC(audio) {
		return true
	}
	return false
}

// Run ...
func Run(args ...string) (string, error) {
	if args == nil {
		args = []string{"-h"}
	}
	cmd := exec.Command("ffprobe", args...)
	cmd.Env = os.Environ()

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		return string(stdout), err
	}

	return string(stdout), nil
}

func filterStream(output string, stream string) string {
	sta := strings.Index(output, stream)
	end := strings.Index(output, "Metadata")
	output = output[sta:end]
	return output
}

// CheckH264 ...
func CheckH264(steam string) bool {
	steam = strings.ToLower(steam)
	if strings.Index(steam, "h264") != -1 {
		return true
	}
	return false
}

// CheckAAC ...
func CheckAAC(steam string) bool {
	steam = strings.ToLower(steam)
	if strings.Index(steam, "aac") != -1 {
		return true
	}
	return false
}
