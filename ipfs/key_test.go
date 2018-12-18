package ipfs

import "testing"

// TestKey_Gen ...
func TestKey_Gen(t *testing.T) {
	config := NewConfig("localhost:5001")
	rlt, err := config.VersionAPI("v0").Key().Gen("tt2", "rsa", 2048)
	t.Log(rlt, err)
}
