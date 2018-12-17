package service

import (
	"github.com/ipfs/go-ipfs-api"
	"strings"
)

var ipfs = InitIPFS("localhost", "5001")

// InitIPFS ...
func InitIPFS(url, port string) *shell.Shell {
	return shell.NewShell(strings.Join([]string{url, port}, ":"))
}
