package ffprobe

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Run ...
func Run(args ...string) ([]byte, error) {
	if args == nil {
		args = []string{"-h"}
	}
	cmd := exec.Command("ffprobe", args...)
	cmd.Env = os.Environ()

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//stderr, err := cmd.StderrPipe()
	//if err != nil {
	//	log.Fatal(err)
	//}
	if err := cmd.Start(); err != nil {
		log.Println(err)
		return nil, err
	}

	//b, err := ioutil.ReadAll(stdout)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//b, err := ioutil.ReadAll(stderr)
	//if err != nil {
	//	log.Fatal(err)
	//}

	if err := cmd.Wait(); err != nil {
		log.Println(err)
		return nil, err
	}

	return stdout, nil
}

// FilterStream ...
func FilterStream(path string) (string, error) {
	b, err := Run(path)
	if err != nil {
		return "", err
	}
	sta := bytes.Index(b, []byte("Stream"))
	end := bytes.Index(b, []byte("Metadata"))
	b = b[sta:end]
	return string(b), nil
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
