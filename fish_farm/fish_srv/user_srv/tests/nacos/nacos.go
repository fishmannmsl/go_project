package main

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"go_project/fish_farm/fish_srv/user_srv/tests/nacos/config"
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
		NamespaceId:         "a3686489-64eb-4d04-87a3-b550520f6db6", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
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
		DataId: "user-srv.json",
		Group:  "dev"})
	if err != nil {
		panic(err)
	}
	//fmt.Println(content) //字符串 - yaml
	serverConfig := config.ServerConfig2{}
	//想要将一个json字符串转换成struct，需要去设置这个struct的tag
	json.Unmarshal([]byte(content), &serverConfig)
	fmt.Println(serverConfig)
	//err = configClient.ListenConfig(vo.ConfigParam{
	//	DataId: "user-web.json",
	//	Group:  "dev",
	//	OnChange: func(namespace, group, dataId, data string) {
	//		fmt.Println("配置文件变化")
	//		fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
	//	},
	//})
	//time.Sleep(3000 * time.Second)

}
