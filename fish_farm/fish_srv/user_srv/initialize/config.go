package initialize

import (
	"encoding/json"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go_project/fish_farm/user_srv/global"
)

/*
读取本地配置文件
*/
func GetEnvInfo(env string) int {
	viper.AutomaticEnv()
	return viper.GetInt(env)
}

//从配置文件中读取对应配置
func InitConfig() {
	var configPath string
	//读取本地环境变量值,设置环境变量后有要重启goland
	debug := GetEnvInfo("ACSvcPort")
	configFilePrefix := "user_srv/config"
	//文件的路径,配置为相对路径
	if debug == 17532 {
		configPath = fmt.Sprintf("%s-debug.yaml", configFilePrefix)
	} else {
		configPath = fmt.Sprintf("%s-pro.yaml", configFilePrefix)
	}

	v := viper.New()
	v.SetConfigFile(configPath)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(&global.NacosConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息: %v", global.NacosConfig)

	//实例化服务端结构体
	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}
	//实例化监听配置信息
	cc := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true, // 在启动的时候不读取缓存在CacheDir的service信息
		LogDir:              "tmp/naocs/log",
		CacheDir:            "tmp/config/cache",
		LogRollingConfig: &constant.ClientLogRollingConfig{
			MaxAge: 3,
		},
		LogLevel: "debug",
	}
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group})
	if err != nil {
		panic(err)
	}
	//将json字符串转化成struct，需要设置struct的tag
	err = json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		zap.S().Fatalf("读取nacos配置失败：%s", err.Error())
	}
	fmt.Println(global.ServerConfig)
}
