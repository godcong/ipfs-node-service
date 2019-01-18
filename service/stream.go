package service

import (
	"github.com/godcong/node-service/openssl"
	"github.com/json-iterator/go"
	"github.com/satori/go.uuid"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// StreamerCallback ...
type StreamerCallback interface {
	Callback(*QueueResult) error
}

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
	FileSource  string
	FileDest    string
	StreamerCallback
}

// NewStreamer ...
func NewStreamer() *Streamer {
	return &Streamer{
		encrypt:     false,
		ID:          uuid.NewV1().String(),
		Key:         "",
		KeyURL:      "",
		KeyName:     "",
		KeyInfoName: "",
		KeyDest:     "",
		//FileName:    util.GenerateRandomString(64),
		FileSource: "",
		FileDest:   "",
	}
}

// NewStreamerWithConfig ...
func NewStreamerWithConfig(cfg *Configure, id string) *Streamer {
	return &Streamer{
		encrypt:     false,
		ID:          DefaultString(id, uuid.NewV1().String()),
		KeyURL:      cfg.Media.KeyURL,
		KeyName:     cfg.Media.KeyFile,
		KeyInfoName: cfg.Media.KeyInfoFile,
		KeyDest:     cfg.Media.KeyDest,
		FileSource:  cfg.Media.Upload,
		FileDest:    cfg.Media.Transfer,
	}
}

// FileName ...
func (s *Streamer) FileName() string {
	_, file := filepath.Split(s.ObjectKey)
	return file
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
	dst := filepath.Join(s.FileDest, s.FileName())
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

// SourceFile ...
func (s *Streamer) SourceFile() string {
	return filepath.Join(s.FileSource, s.ID, s.FileName())
}

// DestPath ...
func (s *Streamer) DestPath() string {
	return filepath.Join(s.FileDest, s.ID)
}

// FromConfig ...
func (s *Streamer) FromConfig(c *Configure) error {
	s.FileDest = c.Media.Transfer
	s.FileSource = c.Media.Upload
	return nil
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
