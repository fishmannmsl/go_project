package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"go_project/fish_farm/user_srv/global"
	"go_project/fish_farm/user_srv/handler"
	"go_project/fish_farm/user_srv/initialize"
	"go_project/fish_farm/user_srv/proto"
	"go_project/fish_farm/user_srv/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//设置IP与端口号
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 50051, "端口号")

	//初始化,注意顺序
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()

	flag.Parse()
	//当部署于生产环境下，将 Port 设置为 0，就可以随机调用端口
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}
	zap.S().Info("ip:", *IP)
	zap.S().Info("port:", *Port)

	server := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{ //这个连接最大的空闲时间，超过就释放，解决proxy等到网络问题（不通知grpc的client和server）
			MaxConnectionIdle: 5 * time.Minute}),
	)
	proto.RegisterUserServer(server, &handler.UserServer{})
	//启动grpc
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	//grpc 注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//注册服务
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", "192.168.0.101", *Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	//生成注册对象,使用 uuid 实现一份代码多个实例
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	registration := api.AgentServiceRegistration{
		Name:    global.ServerConfig.Name,
		ID:      serviceID,
		Port:    *Port,
		Tags:    []string{"srv", "user"},
		Address: "192.168.0.101",
		Check:   check,
	}

	err = client.Agent().ServiceRegister(&registration)
	if err != nil {
		panic(err)
	}

	go func() {
		//使用 gorutine ，否则阻塞无法执行下一步
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	//接受终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")
}
