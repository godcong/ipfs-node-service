package service

import (
	"log"
	"os"
)

// Config ...
type Config struct {
	URL         string //default url
	Upload      string //上传路径
	Transfer    string //转换路径
	M3U8        string //m3u8文件名
	KeyFile     string //key文件名
	KeyInfoFile string //keyFile文件名
}

var config = InitConfig()

// InitConfig ...
func InitConfig() *Config {
	//	TODO:load
	return &Config{
		URL:         "http://localhost:8080",
		Upload:      "upload",
		Transfer:    "transfer",
		M3U8:        "media.m3u8",
		KeyFile:     "key",
		KeyInfoFile: "KeyInfoFile",
	}
}

// Initialize ...
func Initialize(cfg *Config) error {
	if !IsExists(cfg.Upload) {
		err := os.Mkdir(cfg.Upload, os.ModePerm)
		if err != nil {
			return err
		}
	}
	if !IsExists(cfg.Transfer) {
		err := os.Mkdir(cfg.Transfer, os.ModePerm)
		if err != nil {
			return err
		}

	}
	return nil
}

// IsExists ...
func IsExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Panicln(err)
	}
	return true
}
