package ffprobe

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

// Run ...
func Run(args ...string) ([]byte, error) {
	if args == nil {
		args = []string{"-h"}
	}
	cmd := exec.Command("ffprobe", args...)
	cmd.Env = os.Environ()

	//stdout, err := cmd.StdoutPipe()
	//if err != nil {
	//
	//	log.Fatal(err)
	//}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	//b, err := ioutil.ReadAll(stdout)
	//if err != nil {
	//	log.Fatal(err)
	//}

	b, err := ioutil.ReadAll(stderr)
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	return b, nil
}
