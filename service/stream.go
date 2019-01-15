package service

import (
	"github.com/go-redis/redis"
	"github.com/godcong/go-ffmpeg/openssl"
	"github.com/godcong/go-ffmpeg/util"
	"github.com/json-iterator/go"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// Streamer ...
type Streamer struct {
	encrypt     bool
	ID          string
	Key         string
	ObjectKey   string
	KeyURL      string
	KeyName     string
	KeyInfoName string
	KeyDest     string
	FileName    string
	FileSource  string
	FileDest    string
}

// NewStreamer ...
func NewStreamer(id string) *Streamer {
	return &Streamer{
		encrypt:     false,
		ID:          id,
		Key:         "",
		KeyURL:      "",
		KeyName:     "",
		KeyInfoName: "",
		KeyDest:     "",
		FileName:    util.GenerateRandomString(64),
		FileSource:  "",
		FileDest:    "",
	}
}

// Encrypt ...
func (s *Streamer) Encrypt() bool {
	return s.encrypt
}

// SetEncrypt ...
func (s *Streamer) SetEncrypt(encrypt bool) {
	s.encrypt = true
	s.KeyURL = config.Media.KeyURL
	s.KeyName = config.Media.KeyFile
	s.KeyInfoName = config.Media.KeyInfoFile
	s.KeyDest = config.Media.KeyDest
}

// KeyFile ...
func (s *Streamer) KeyFile() string {
	var err error
	dst := filepath.Join(s.FileDest, s.FileName)
	err = os.Mkdir(dst, os.ModePerm)
	if err != nil {
		log.Println(err)
		return ""
	}

	err = openssl.KeyFile(s.KeyDest, s.KeyName, s.Key, s.KeyInfoName, s.KeyURL, true)
	if err != nil {
		log.Println(err)
		return ""
	}

	return dst + "/" + s.KeyInfoName
}

// JSON ...
func (s *Streamer) JSON() string {
	st, err := jsoniter.MarshalToString(s)
	if err != nil {
		log.Println(err)
		return ""
	}
	return st
}

// ParseStreamer ...
func ParseStreamer(s string) *Streamer {
	var st Streamer
	err := jsoniter.UnmarshalFromString(s, &st)
	if err != nil {
		return nil
	}
	return &st
}

// StreamQueue ...
type StreamQueue struct {
	infos []*Streamer
	lock  sync.RWMutex
}

// NewRedisQueue ...
func NewRedisQueue() *redis.Client {
	return newRedisWithDB(RedisQueueIndex)
}

// NewStreamQueue ...
func NewStreamQueue() *StreamQueue {
	return &StreamQueue{
		infos: []*Streamer{},
	}
}

// Push ...
func (s *StreamQueue) Push(info *Streamer) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.infos = append(s.infos, info)
}

// Pop ...
func (s *StreamQueue) Pop() *Streamer {
	s.lock.Lock()
	defer s.lock.Unlock()
	info := s.infos[0]
	s.infos = s.infos[1:len(s.infos)]

	return info
}

// Front ...
func (s *StreamQueue) Front() *Streamer {
	s.lock.RLock()
	defer s.lock.RUnlock()
	info := s.infos[0]

	return info
}

// IsEmpty ...
func (s *StreamQueue) IsEmpty() bool {
	return len(s.infos) == 0
}

// Size ...
func (s *StreamQueue) Size() int {
	return len(s.infos)
}

// Clear ...
func (s *StreamQueue) Clear() {
	s.infos = []*Streamer{}
}
