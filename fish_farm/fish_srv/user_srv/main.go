package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"net"

	"go_project/fish_farm/user_srv/handler"
	"go_project/fish_farm/user_srv/proto"
)

func main() {
	//设置IP与端口号
	IP := flag.String("ip","0.0.0.0","ip地址")
	Port := flag.Int("port",9080,"端口号")
	flag.Parse()
	fmt.Println("ip:",*IP)
	fmt.Println("port:",*Port)

	server := grpc.NewServer()
	proto.RegisterUserServer(server,&handler.UserServer{})
	//启动grpc
	lis,err := net.Listen("tcp",fmt.Sprintf("%s:%d",*IP,*Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	err = server.Serve(lis)
	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}
}
