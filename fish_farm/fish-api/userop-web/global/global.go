package global

import (
	ut "github.com/go-playground/universal-translator"
	
	"go_project/fish_farm/fish-api/userop-web/config"
	"go_project/fish_farm/fish-api/userop-web/proto"
)

var (
	Trans ut.Translator

	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	NacosConfig *config.NacosConfig = &config.NacosConfig{}

	GoodsSrvClient proto.GoodsClient

	MessageClient proto.MessageClient
	AddressClient proto.AddressClient
	UserFavClient proto.UserFavClient
)
