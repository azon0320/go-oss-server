package config

import (
	"fmt"
	"github.com/dormao/go-oss-server/internal/context/consts"
	"github.com/jinzhu/configor"
	"os"
)

var Config = struct {
	AccessKey    string
	AccessSecret string
	Bucket       string
	OutputMode   string
	ServeAddress string
	BaseURL      string
	StorePath    string
	Provider     struct {
		Type          string
		FilePath      string
		PostgresURI   string
		DbAutoMigrate bool
	}
}{}

func init() {
	err := configor.Load(&Config, os.Getenv(consts.EnvConfigFile))
	if err != nil {
		fmt.Errorf("error while loading oss config file: %s", err)
	}
}
