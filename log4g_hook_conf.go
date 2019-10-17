package log4g

import (
	"github.com/BurntSushi/toml"
)

var configPath = "./log4g.toml"

func InitAppenders() (*Appenders, error) {
	var appenders Appenders
	if _, err := toml.DecodeFile(configPath, &appenders); nil != err {
		return nil, err
	}
	return &appenders, nil
}

type Appenders struct {
	Appender []Appender
}

type Appender struct {
	Filename   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	LocalTime  bool
	Compress   bool
	MaxLevel   string
	MinLevel   string
}
