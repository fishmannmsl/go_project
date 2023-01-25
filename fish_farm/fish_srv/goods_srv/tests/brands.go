package main

import (
	"context"
	"fmt"
	"go_project/fish_farm/fish_srv/goods_srv/proto"
	"google.golang.org/grpc"
)

/*
测试功能接口文件
*/
var (
	BrandClient proto.GoodsClient
	conn        *grpc.ClientConn
)

/*
查询用户方法测试
*/
func TestGetBrandList() {
	rsp, err := BrandClient.BrandList(context.Background(), &proto.BrandFilterRequest{})
	fmt.Println(rsp)
	if err != nil {
		panic(err)
	}
	for _, brand := range rsp.Data {
		fmt.Println(brand.Name)
	}
}
