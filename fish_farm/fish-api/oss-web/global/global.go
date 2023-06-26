package global

import (
	ut "github.com/go-playground/universal-translator"
	"go_project/fish_farm/fish-api/oss-web/config"
)

var (
	Trans ut.Translator

	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	NacosConfig *config.NacosConfig = &config.NacosConfig{}
)
