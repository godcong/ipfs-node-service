package oss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/godcong/ipfs-node-service/config"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

// OSS ...
type OSS struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
	*oss.Bucket
}

// BucketServer ...
type BucketServer struct {
	config       *config.Configure
	server       map[string]*OSS
	current      *OSS
	DownloadPath string
	UploadPath   string
	PartSize     int64
	Routines     oss.Option
	Checkpoint   oss.Option
	Progress     oss.Option
}

var server *BucketServer

// NewOSS ...
func NewOSS(oss *config.OSS) *OSS {
	return &OSS{
		Endpoint:        config.DefaultString(oss.Endpoint, ""),
		AccessKeyID:     config.DefaultString(oss.AccessKeyID, ""),
		AccessKeySecret: config.DefaultString(oss.AccessKeySecret, ""),
		BucketName:      config.DefaultString(oss.BucketName, ""),
	}
}

// Connect ...
func (o *OSS) Connect() error {
	client, err := oss.New(o.Endpoint, o.AccessKeyID, o.AccessKeySecret)
	if err != nil {
		return fmt.Errorf("failed to create new client: %s", err)
	}

	bucket, err := client.Bucket(o.BucketName)
	if err != nil {
		return fmt.Errorf("failed to get bucket: %s", err)
	}
	o.Bucket = bucket
	return nil
}

// NewBucketServer ...
func NewBucketServer(cfg *config.Configure) *BucketServer {
	var s BucketServer
	s.config = cfg
	s.server = make(map[string]*OSS)
	s.DownloadPath = config.DefaultString(cfg.Media.Download, "download")
	s.UploadPath = config.DefaultString(cfg.Media.Upload, "upload")
	s.PartSize = 100 * 1024 * 1024
	s.Routines = oss.Routines(5)
	s.Checkpoint = oss.Checkpoint(true, "./cp")
	s.Progress = oss.Progress(&s)

	for _, val := range cfg.OSS {
		oss := NewOSS(&val)
		err := oss.Connect()
		if err != nil {
			log.Panic(err)
			return nil
		}
		s.server[oss.BucketName] = oss
		s.current = oss
	}
	return &s
}

// InitOSS ...
func InitOSS(cfg *config.Configure) {
	server = NewBucketServer(cfg)
}

// 定义进度变更事件处理函数。
func (s *BucketServer) ProgressChanged(event *oss.ProgressEvent) {
	switch event.EventType {
	case oss.TransferStartedEvent:
		log.Printf("Transfer Started, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case oss.TransferDataEvent:
		log.Printf("\rTransfer Data, ConsumedBytes: %d, TotalBytes %d, %d%%.",
			event.ConsumedBytes, event.TotalBytes, event.ConsumedBytes*100/event.TotalBytes)
	case oss.TransferCompletedEvent:
		log.Printf("\nTransfer Completed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)

	case oss.TransferFailedEvent:
		log.Printf("\nTransfer Failed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	default:
	}
}

func (s *BucketServer) Current() *OSS {
	var b bool
	if s.current == nil {
		s.current, b = s.Server("")
		if !b {
			log.Panic("no oss founded")
		}
	}
	return s.current
}

// Server ...
func (s *BucketServer) Server(name string) (*OSS, bool) {
	if name == "" {
		for _, v := range s.server {
			if v != nil {
				return v, true
			}
		}
	}
	v, b := s.server[name]
	return v, b
}

// Node ...
type Node interface {
	SetBucketName(string)
	BucketName() string
	ObjectKey() string
	OutputName() string
}

type data struct {
	bucketName string
	objectKey  string
	outputName string
}

func (d *data) SetBucketName(name string) {
	d.bucketName = name
}

func (d *data) BucketName() string {
	return d.bucketName
}

func (d *data) Options() []oss.Option {
	panic("implement me")
}

// Path ...
func (d *data) OutputName() string {
	return d.outputName
}

// NewNode ...
func NewNode(objectKey string, outputName string) Node {
	return &data{
		objectKey:  objectKey,
		outputName: outputName,
	}
}

// ObjectKey ...
func (d *data) ObjectKey() string {
	return d.objectKey
}

// SetObjectKey ...
func (d *data) SetObjectKey(objectKey string) {
	d.objectKey = objectKey
}

// FileName ...
func FileName(objectKey string) string {
	_, file := filepath.Split(objectKey)
	return file
}

// Download ...
func (s *BucketServer) Download(p Node) error {
	fp := filepath.Join(s.DownloadPath, p.OutputName())
	abs, e := filepath.Abs(s.DownloadPath)
	if e != nil {
		return e
	}
	_ = os.MkdirAll(abs, os.ModePerm)
	e = s.Current().DownloadFile(p.ObjectKey(), fp, s.PartSize, s.Routines, s.Progress, s.Checkpoint)
	if e != nil {
		return e
	}
	return nil
}

// Upload ...
func (s *BucketServer) Upload(p Node) error {
	fp := filepath.Join(s.DownloadPath, p.OutputName())
	//abs, e := filepath.Abs(s.DownloadPath)
	//if e != nil {
	//	return e
	//}
	//_ = os.MkdirAll(abs, os.ModePerm)
	err := s.Current().UploadFile(p.ObjectKey(), fp, s.PartSize, s.Routines, s.Progress, s.Checkpoint)
	if err != nil {
		return err
	}
	return nil
}

// URL ...
func (s *BucketServer) URL(p Node) (string, error) {
	signedURL, err := s.Current().SignURL(p.ObjectKey(), oss.HTTPGet, 60*60*24)
	if err != nil {
		return "", err
	}
	return signedURL, err

}

// IsExist ...
func (o *OSS) IsExist(p Node) bool {
	exist, err := o.Bucket.IsObjectExist(p.ObjectKey())
	if err != nil {
		log.Println(err)
		return false
	}
	return exist
}

// Server ...
func Server() *BucketServer {
	return server
}
