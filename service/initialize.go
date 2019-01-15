package service

import (
	"github.com/pelletier/go-toml"
	"log"
	"os"
)

// Media ...
type Media struct {
	Upload      string `toml:"upload"`        //上传路径
	Transfer    string `toml:"transfer"`      //转换路径
	M3U8        string `toml:"m3u8"`          //m3u8文件名
	KeyURL      string `toml:"key_url"`       //default url
	KeyDest     string `toml:"key_dest"`      //key 文件输出目录
	KeyFile     string `toml:"key_file"`      //key文件名
	KeyInfoFile string `toml:"key_info_file"` //keyFile文件名
}

// IPFS ...
type IPFS struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
}

// GRPC ...
type GRPC struct {
	Enable bool `toml:"enable"`
}

// REST ...
type REST struct {
	Enable bool `toml:"enable"`
}

// Queue ...
type Queue struct {
	Addr     string
	Password string
}

// Configure ...
type Configure struct {
	Media Media `toml:"media"`
	GRPC  GRPC  `toml:"grpc"`
	REST  REST  `toml:"rest"`
	IPFS  IPFS  `toml:"ipfs"`
}

var config *Configure

// Initialize ...
func Initialize(filePath ...string) error {
	if filePath == nil {
		filePath = []string{"config.toml"}
	}

	cfg := LoadConfig(filePath[0])

	if !IsExists(cfg.Media.Upload) {
		err := os.Mkdir(cfg.Media.Upload, os.ModePerm)
		if err != nil {
			return err
		}
	}
	if !IsExists(cfg.Media.Transfer) {
		err := os.Mkdir(cfg.Media.Transfer, os.ModePerm)
		if err != nil {
			return err
		}
	}

	if !IsExists(cfg.Media.KeyDest) {
		err := os.Mkdir(cfg.Media.KeyDest, os.ModePerm)
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
