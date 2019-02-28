package service

import (
	"github.com/godcong/ipfs-media-service/util"
	log "github.com/sirupsen/logrus"
	"testing"
)

// TestRedisHSET ...
func TestRedisHSET(t *testing.T) {

}

// TestNewRedisQueue ...
func TestNewRedisQueue(t *testing.T) {
	for i := 0; i < 100; i++ {
		s := NewStreamer()
		s.ObjectKey = util.GenerateRandomString(64)
		Push(s)
	}
	count := 0
	for v := Pop(); v != nil; v = Pop() {
		log.Println(v)
		count++
	}

	log.Println(count)

}
