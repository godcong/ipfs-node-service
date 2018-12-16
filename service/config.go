package service

type Config struct {
	Upload   string //上传路径
	Transfer string //转换路径
	Key      string //key文件名
	KeyFile  string //keyFile文件名
}

var config = InitConfig()

func InitConfig() *Config {
	//	TODO:load
	return &Config{
		Upload:   "./upload/",
		Transfer: "./transfer/",
		Key:      "key",
		KeyFile:  "KeyFile",
	}
}
