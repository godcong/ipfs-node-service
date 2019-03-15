package service

import (
	"github.com/godcong/ipfs-node-service/config"
	"github.com/godcong/ipfs-node-service/openssl"
	"github.com/google/uuid"
	"github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
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
	config      *config.Configure
	Encrypt     bool
	ID          string
	Key         string
	ObjectKey   string
	KeyURL      string
	KeyName     string
	KeyInfoName string
	KeyDest     string
	Download    string
	Transfer    string
	Callback    string
}

// NewStreamer ...
func NewStreamer() *Streamer {
	return &Streamer{
		Encrypt:     false,
		ID:          uuid.New().String(),
		Key:         "",
		KeyURL:      "",
		KeyName:     "",
		KeyInfoName: "",
		KeyDest:     "",
		//FileName:    util.GenerateRandomString(64),
		Download: "",
		Transfer: "",
	}
}

// NewStreamerWithConfig ...
func NewStreamerWithConfig(cfg *config.Configure, id string) *Streamer {
	return &Streamer{
		config:      cfg,
		Encrypt:     false,
		ID:          config.DefaultString(id, uuid.New().String()),
		Key:         "",
		ObjectKey:   "",
		KeyURL:      config.DefaultString(cfg.Media.KeyURL, "http://localhost:8080"),
		KeyName:     config.DefaultString(cfg.Media.KeyFile, "key"),
		KeyInfoName: config.DefaultString(cfg.Media.KeyInfoFile, "KeyInfoFile"),
		KeyDest:     config.DefaultString(cfg.Media.KeyDest, "output_key"),
		Download:    config.DefaultString(cfg.Media.Download, "download"),
		Transfer:    config.DefaultString(cfg.Media.Transfer, "transfer"),
		Callback:    "",
	}
}

// FileName ...
func (s *Streamer) FileName() string {
	_, file := filepath.Split(s.ObjectKey)
	return file
}

// SetEncrypt ...
func (s *Streamer) SetEncrypt(keyURL string) {
	if keyURL != "" {
		s.Encrypt = true
		s.KeyURL = config.DefaultString(keyURL, s.KeyURL)
	}
}

// KeyFile ...
func (s *Streamer) KeyFile() string {
	var err error
	dst := filepath.Join(s.KeyDest, s.ID)
	abs, err := filepath.Abs(dst)
	if err != nil {
		log.Error(err)
		return ""
	}
	err = os.Mkdir(abs, os.ModePerm)
	if err != nil {
		log.Error(err)
		return ""
	}

	err = openssl.KeyFile(abs, s.KeyName, s.Key, s.KeyInfoName, s.KeyURL, true)
	if err != nil {
		log.Error(err)
		return ""
	}

	return filepath.Join(abs, s.KeyInfoName)
}

// SourceFile ...
func (s *Streamer) SourceFile() string {
	return filepath.Join(s.Download, s.ID)
}

// DestPath ...
func (s *Streamer) DestPath() string {
	return filepath.Join(s.Transfer, s.ID)
}

// FromConfig ...
func (s *Streamer) FromConfig(cfg *config.Configure) error {
	s.Transfer = cfg.Media.Transfer
	s.Download = cfg.Media.Download
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
