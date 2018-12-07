package ffmpeg

import "testing"

func TestRun(t *testing.T) {
	b, err := Run("-version")
	t.Log(string(b), err)
}
