package global

import (
	"go_project/fish_farm/fish_srv/user_srv/config"
	"gorm.io/gorm"
)

//将DB作为全局变量
var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  *config.NacosConfig
)
