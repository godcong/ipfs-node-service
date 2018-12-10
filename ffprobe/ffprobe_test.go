package ffprobe

import "testing"

// TestRun ...
func TestRun(t *testing.T) {
	run, e := Run("../tmp/4ltifGZK4mfkK5EEbCBEhnhv8puSYU890dq34y5sRXXMs6k44Zjm87BhIplDwoby")
	t.Log(string(run), e)
}
