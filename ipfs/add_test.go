package ipfs

import "testing"

func TestApi_AddDir(t *testing.T) {
	cfg := NewConfig("localhost:5001")
	api := cfg.VersionAPI("v0")
	dir, e := api.AddDir("D:\\workspace\\goproject\\ipfs-node-service\\transfer\\e504cbc9-d355-486d-abf1-246607f89a3b\\")
	t.Log(dir, e)
}
