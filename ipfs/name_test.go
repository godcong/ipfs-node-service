package ipfs

import "testing"

// TestName_Publish ...
func TestName_Publish(t *testing.T) {
	config := NewConfig("localhost:5001")
	rlt, err := config.VersionAPI("v0").Name().PublishD("/ipfs/QmeKefn6f3yh1A6J6QvbidQEwAgJbg8VD8UDeFdM75JjrN")
	t.Log(rlt, err)
}
