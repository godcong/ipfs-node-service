package service

import (
	"github.com/pelletier/go-toml"
	"log"
	"os"
)

// IPFS ...
type IPFS struct {
	Upload      string `toml:"upload"`        //上传路径
	Transfer    string `toml:"transfer"`      //转换路径
	M3U8        string `toml:"m3u8"`          //m3u8文件名
	KeyURL      string `toml:"key_url"`       //default url
	KeyFile     string `toml:"key_file"`      //key文件名
	KeyInfoFile string `toml:"key_info_file"` //keyFile文件名
}

// GRPC ...
type GRPC struct {
	Enable string `toml:"enable"`
}

// REST ...
type REST struct {
	Enable string `toml:"enable"`
}

// Configure ...
type Configure struct {
	IPFS IPFS `toml:"ipfs"`
	GRPC GRPC `toml:"grpc"`
	REST REST `toml:"rest"`
}

var config *Configure

// Initialize ...
func Initialize(filePath ...string) error {
	if filePath == nil {
		filePath = []string{"config.toml"}
	}

	cfg := LoadConfig(filePath[0])

	if !IsExists(cfg.IPFS.Upload) {
		err := os.Mkdir(cfg.IPFS.Upload, os.ModePerm)
		if err != nil {
			return err
		}
	}
	if !IsExists(cfg.IPFS.Transfer) {
		err := os.Mkdir(cfg.IPFS.Transfer, os.ModePerm)
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
func LoadConfig(filePath string) *Configure {
	var cfg Configure
	openFile, err := os.OpenFile(filePath, os.O_RDONLY|os.O_SYNC, os.ModePerm)
	if err != nil {
		panic(err.Error())
	}
	decoder := toml.NewDecoder(openFile)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err.Error())
	}
	log.Printf("config: %+v", cfg)
	return &cfg
}

// Config ...
func Config() *Configure {
	if config == nil {
		panic("nil config")
	}
	return config
}
