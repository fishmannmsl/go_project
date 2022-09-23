package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"strings"
	"time"

	"github.com/anaskhan96/go-password-encoder"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"

	"go_project/fish_farm/user_srv/global"
	"go_project/fish_farm/user_srv/model"
	"go_project/fish_farm/user_srv/proto"
)

/*由于版本更新,proto会生成一个私有mustEmbedUnimplementedUserServer()函数
导致实现这个 proto.UserServer{} 接口时无法实现，
所以我们可以通过 在UserServer 结构体内添加 proto.UnimplementedUserServer
来实现 proto.UserServer{} 这个接口
*/
type UserServer struct {
	proto.UnimplementedUserServer
}

//将go语言中的用户结构转化为proto中的message
func ModelToResponse(user model.User) proto.UserInfoResponse{
	//在grpc的message中有自己的默认值，不能将赋值nil，容易出错，所以将不为 nil 的属性加入
	userInfoRsp := proto.UserInfoResponse{
		Id: user.ID,
		PassWord: user.Password,
		Mobile: user.Mobile,
		NickName: user.NickName,
		Gender:user.Gender,
		Role:uint32(user.Role),
	}
	if user.Birthday != nil {
		userInfoRsp.Birthday = uint64(user.Birthday.Unix())
	}
	return userInfoRsp
}

//分页
func Paginate(page,pageSize int) func(db *gorm.DB) *gorm.DB {
	return func (db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

//获取用户列表
func (u UserServer) GetUserList(ctx context.Context,req *proto.PageInfo) (*proto.UserListResponse, error) {
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)
	global.DB.Scopes(Paginate(int(req.Pn),int(req.PSzie))).Find(&users)

	for _,user := range users{
		userInfoRsp := ModelToResponse(user)
		rsp.Data = append(rsp.Data,&userInfoRsp)
	}
	fmt.Println(rsp.Data)

	return rsp, nil
}

//通过手机查询用户
func (u UserServer) GetUserByMobile(ctx context.Context,req *proto.MobileRequest) (*proto.UserInfoResponse, error)  {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	//用户不存在
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound,"用户不存在")
	}
	//其他错误
	if result.Error != nil {
		return nil, result.Error
	}

	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp,nil
}

//通过ID查询用户
func (u UserServer) GetUserByID(ctx context.Context,req *proto.IDRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user,req.Id)
	//用户不存在
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound,"用户不存在")
	}
	//其他错误
	if result.Error != nil {
		return nil, result.Error
	}

	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp,nil
}

//创建用户
func (u UserServer) CreateUser(ctx context.Context,req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	var user model.User
	//判断手机是否已被注册
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 1 {
		return nil,status.Errorf(codes.AlreadyExists,"该手机号已被注册")
	}
	user.Mobile = req.Mobile
	user.NickName = req.NickName

	//密码加密
	options := &password.Options{16, 100, 32, sha512.New}//md5.New可改为sha512.New()
	salt, encodedPwd := password.Encode(req.PassWord, options)
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s",salt,encodedPwd)

	result = global.DB.Create(&user)
	if result.Error != nil {
		//如果出错只能是内部的错误
		return nil, status.Errorf(codes.Internal,result.Error.Error())
	}

	userInfo := ModelToResponse(user)
	return &userInfo,nil
}

//更新用户信息
func (u UserServer) UpdateUser(ctx context.Context,req *proto.UpdateUserInfo) (*emptypb.Empty, error) {
	//用户必须存在
	var user model.User
	result := global.DB.First(&user,req.Id)
	if result.RowsAffected == 0 {
		return nil,status.Errorf(codes.NotFound,"用户不存在")
	}

	birthDay := time.Unix(int64(req.Birthday),0)
	user.NickName = req.NickName
	user.Birthday = &birthDay
	user.Gender = req.Gender

	//更新数据
	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil,status.Errorf(codes.Internal,result.Error.Error())
	}
	return &empty.Empty{},nil
}

//校验密码
func (u UserServer) CheckPassWord(ctx context.Context,req *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	//解析密码并验证
	options := &password.Options{16, 100, 32, sha512.New}
	passwordInfo := strings.Split(req.EncryptedPassword,"$")
	check := password.Verify(req.PassWord, passwordInfo[2], passwordInfo[3], options)
	return &proto.CheckResponse{Success:check},nil
}