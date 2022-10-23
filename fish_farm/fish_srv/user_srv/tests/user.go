package main

import (
	"context"
	"fmt"
	"go_project/fish_farm/user_srv/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/*
测试功能接口文件
*/
var (
	userClient proto.UserClient
	conn       *grpc.ClientConn
)

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	userClient = proto.NewUserClient(conn)
}

/*
查询用户方法测试
*/
func TestGetUserList() {
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSzie: 2,
	})
	fmt.Println(rsp)
	if err != nil {
		panic(err)
	}
	for _, user := range rsp.Data {
		fmt.Println(user.Mobile, user.NickName, user.PassWord)
		CheckRsp, err := userClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			PassWord:          "admin123",
			EncryptedPassword: user.PassWord,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(CheckRsp.Success)
	}
}

/*
创建用户方法测试
*/
func TestCreateUser() {
	for i := 0; i < 10; i++ {
		rsp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
			NickName: fmt.Sprintf("zgnbs%d", i),
			Mobile:   fmt.Sprintf("1276165457%d", i),
			PassWord: "admin123",
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(rsp.Id)
	}
}

/*
更新用户方法测试
以及通过ID进行查询
*/
func TestUpdateUser() {
	_, err := userClient.UpdateUser(context.Background(), &proto.UpdateUserInfo{
		Id:       int32(1),
		NickName: "dsx",
		Gender:   "male",
	})
	if err != nil {
		panic(err)
	}
	rsp, err := userClient.GetUserByID(context.Background(), &proto.IDRequest{
		Id: int32(1),
	})
	fmt.Println(rsp.Id, rsp.NickName, rsp.Gender)

}

func main() {
	Init()
	TestGetUserList()
	//TestCreateUser()
	//TestUpdateUser()
	conn.Close()
}
