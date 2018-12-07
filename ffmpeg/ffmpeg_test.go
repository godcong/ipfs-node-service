package openssl

import "testing"

// TestRun ...
func TestRun(t *testing.T) {
	b, err := Run("-version")
	t.Log(string(b), err)
}
