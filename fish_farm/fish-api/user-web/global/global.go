package global

import (
	ut "github.com/go-playground/universal-translator"
	"go_project/fish_farm/fish-api/user-web/config"
	"go_project/fish_farm/fish-api/user-web/proto"
)

/*
放置全局变量
*/
var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	NacosConfig  *config.NacosConfig  = &config.NacosConfig{}

	Trans ut.Translator

	UserSrvClient proto.UserClient
)
