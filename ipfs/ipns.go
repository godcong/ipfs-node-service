package ipfs

import (
	"os"
	"os/exec"
)

// Run ...
func Run(args ...string) ([]byte, error) {
	if args == nil {
		args = []string{"help"}
	}
	cmd := exec.Command("ipfs", args...)
	cmd.Env = os.Environ()
	cmd.Env = os.Environ()

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return stdout, err
	}

	return stdout, nil
}

// KeyGen ...
func KeyGen(name string) (string, error) {
	b, err := Run("key", "gen", "-t", "rsa", "-s", "2048", name)
	return string(b), err
}
