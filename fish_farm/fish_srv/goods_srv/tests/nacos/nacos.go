package main

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"go_project/fish_farm/fish_srv/goods_srv/tests/nacos/config"
)

func main() {
	//实例化服务端结构体
	sc := []constant.ServerConfig{
		{
			IpAddr: "192.168.0.109",
			Port:   8848,
		},
	}
	cc := constant.ClientConfig{
		NamespaceId:         "6195d0cc-4563-49ba-bf64-ea454db5d2b6", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/config/log",
		CacheDir:            "tmp/config/cache",
		LogLevel:            "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "goods-srv.json",
		Group:  "dev"})
	if err != nil {
		panic(err)
	}
	serverConfig := config.ServerConfig2{}
	//想要将一个json字符串转换成struct，需要去设置这个struct的tag
	json.Unmarshal([]byte(content), &serverConfig)
	fmt.Println(serverConfig)
}
