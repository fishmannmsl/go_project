package initialize

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go_project/fish_farm/fish-api/goods-web/global"
	"go_project/fish_farm/fish-api/goods-web/proto"
)

/*
初始化 Grpc 服务
*/
func InitSrvConn() {
	//从consul内获取 user_srv 服务
	userconn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s",
			global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port,
			global.ServerConfig.UserInfo.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务器失败】")
	}
	//事先创建好 GRPC 连接就不用再次进行 tcp 的三次握手协议，提高效率
	//一个连接多个 groutine 共用，性能降低，使用连接池来改善
	userSrvClient := proto.NewUserClient(userconn)
	global.UserSrvClient = userSrvClient
}

/*
 * @ v1.0
 * @ 未使用github.com/mbobakov/grpc-consul-resolver
 */
func InitSrvConn2() {
	//从注册中心获取到用户服务的信息
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	userSrvHost := ""
	userSrvPort := 0
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"",
		global.ServerConfig.UserInfo.Name))
	if err != nil {
		panic(err)
	}
	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}
	if userSrvHost == "" {
		zap.S().Fatal("[InitSrvConn] 连接【用户服务】失败")
		return
	}
	//拨号连接用户grpc服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost,
		userSrvPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务器失败】",
			"msg", err.Error())
	}
	//事先创建好 GRPC 连接就不用再次进行 tcp 的三次握手协议，提高效率
	//一个连接多个 groutine 共用，性能降低，使用连接池来改善
	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient

}
