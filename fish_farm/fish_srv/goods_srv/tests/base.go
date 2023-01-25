package main

import (
	"go_project/fish_farm/fish_srv/goods_srv/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	BrandClient = proto.NewGoodsClient(conn)
}

func main() {
	Init()
	TestGetBrandList()
	conn.Close()
}
