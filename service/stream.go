package service

import (
	"github.com/godcong/go-ffmpeg/openssl"
	"log"
	"os"
	"sync"
)

// StreamInfo ...
type StreamInfo struct {
	encrypt  bool
	key      string
	fileName string
	uri      string
	src      string
	dst      string
}

// NewStreamer ...
func NewStreamer(key string, fileName string) *StreamInfo {
	return &StreamInfo{key: key, fileName: fileName}
}

// Encrypt ...
func (s *StreamInfo) Encrypt() bool {
	return s.encrypt
}

// SetEncrypt ...
func (s *StreamInfo) SetEncrypt(encrypt bool) {
	s.encrypt = encrypt
}

// Dst ...
func (s *StreamInfo) Dst() string {
	return s.dst
}

// SetDst ...
func (s *StreamInfo) SetDst(dst string) {
	s.dst = dst
}

// Src ...
func (s *StreamInfo) Src() string {
	return s.src
}

// SetSrc ...
func (s *StreamInfo) SetSrc(src string) {
	s.src = src
}

// FileName ...
func (s *StreamInfo) FileName() string {
	return s.fileName
}

// SetFileName ...
func (s *StreamInfo) SetFileName(fileName string) {
	s.fileName = fileName
}

// Key ...
func (s *StreamInfo) Key() string {
	return s.key
}

// SetKey ...
func (s *StreamInfo) SetKey(key string) {
	s.key = key
}

// URI ...
func (s *StreamInfo) URI() string {
	return s.uri
}

// SetURI ...
func (s *StreamInfo) SetURI(uri string) {
	s.uri = uri
}

// KeyFile ...
func (s *StreamInfo) KeyFile() string {
	var err error
	dst := s.dst + s.fileName
	err = os.Mkdir(dst, os.ModePerm)
	if err != nil {
		log.Println(err)
		return ""
	}

	err = openssl.KeyFile(dst, config.KeyFile, s.key, config.KeyInfoFile, s.uri, true)
	if err != nil {
		log.Println(err)
		return ""
	}

	return dst + "/" + config.KeyInfoFile
}

// StreamQueue ...
type StreamQueue struct {
	infos []*StreamInfo
	lock  sync.RWMutex
}

// NewStreamQueue ...
func NewStreamQueue() *StreamQueue {
	return &StreamQueue{
		infos: []*StreamInfo{},
	}
}

// Push ...
func (s *StreamQueue) Push(info *StreamInfo) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.infos = append(s.infos, info)
}

// Pop ...
func (s *StreamQueue) Pop() *StreamInfo {
	s.lock.Lock()
	defer s.lock.Unlock()
	info := s.infos[0]
	s.infos = s.infos[1:len(s.infos)]

	return info
}

// Front ...
func (s *StreamQueue) Front() *StreamInfo {
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
	s.infos = []*StreamInfo{}
}
