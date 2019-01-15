package service

import (
	"github.com/pelletier/go-toml"
	"log"
	"os"
)

// Config ...
type Config struct {
	Upload      string `json:"upload"`        //上传路径
	Transfer    string `json:"transfer"`      //转换路径
	M3U8        string `json:"m3u8"`          //m3u8文件名
	KeyURL      string `json:"key_url"`       //default url
	KeyFile     string `json:"key_file"`      //key文件名
	KeyInfoFile string `json:"key_info_file"` //keyFile文件名
}

var config *Config

// Initialize ...
func Initialize(filePath ...string) error {
	if filePath == nil {
		filePath = []string{"config.toml"}
	}

	cfg := LoadConfig(filePath[0])

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

	config = cfg

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

// LoadConfig ...
func LoadConfig(filePath string) *Config {
	var cfg Config
	openFile, err := os.OpenFile(filePath, os.O_RDONLY|os.O_SYNC, os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
	decoder := toml.NewDecoder(openFile)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err.Error())
	}
	return &cfg
}
