package api

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt/v4"
	"go_project/fish_farm/fish-api/user-web/middlewares"
	"go_project/fish_farm/fish-api/user-web/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go_project/fish_farm/fish-api/user-web/forms"
	"go_project/fish_farm/fish-api/user-web/global"
	"go_project/fish_farm/fish-api/user-web/global/response"
	"go_project/fish_farm/fish-api/user-web/proto"
)

/*
将错误格式化
*/
func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

/*
将Grpc的错误(code) 转化为 Http的错误状态码
*/
func HandleGrpcErrorToHttp(err error, ctx *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound: // 404
				ctx.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal: //内部错误，不能暴露给用户，返回 500
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument: //参数错误
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable: //数据损坏或不可用
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": e.Code(),
				})
			}
			return
		}
	}

}

/*
调用验证器，处理错误信息,将grpc中的错误转化为http的状态码
*/
func HandleValidatorError(ctx *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
}

/*
获取用户列表
*/
func GetUserList(ctx *gin.Context) {
	//获取登录用户的ID
	claims, _ := ctx.Get("claims")
	//将取得的数据转化为 models.CustomClaims 结构体类型
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("访问用户: %d", currentUser.ID)
	//生成grpc的client并调用接口
	//userSrvClient := proto.NewUserClient(userConn)

	pn := ctx.DefaultQuery("pn", "0")
	pInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10")
	PSizeInt, _ := strconv.Atoi(pSize)
	rsp, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pInt),
		PSzie: uint32(PSizeInt),
	})
	if err != nil {
		//不能随便抛出异常，导致程序被中断
		zap.S().Errorw("[GetUserList] 查询 【用户列表】 失败")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		//data := make(map[string]interface{})

		user := response.UserResponse{
			Id:       value.Id,
			PassWord: value.PassWord,
			Mobile:   value.Mobile,
			NickName: value.NickName,
			Birthday: response.JsonTime(time.Unix(int64(value.Birthday), 0)), //获取 time.Time 类型当前时间
			//Birthday: time.Time(time.Unix(int64(value.Birthday),0)).Format("1970-01-01"),
			Gender: value.Gender,
			Role:   value.Role,
		}
		//在python内常用，但是在go内，优先采用结构体，为了简便也可以采用这种方式
		//data["id"] = value.Id
		//data["name"] = value.NickName
		//data["birthday"] = value.Birthday
		//data["gender"] = value.Gender
		//data["mobile"] = value.Mobile

		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)
}

/*
密码登录
*/
func PassWordLogin(ctx *gin.Context) {
	//表单验证，格式是否正确，使用了第三方包 validator ,内置验证规则(正则)
	passwordLoginForms := forms.PassWordLoginForms{}
	if err := ctx.ShouldBindJSON(&passwordLoginForms); err != nil {
		HandleValidatorError(ctx, err)
		return
	}
	//对输入的图片验证码进行验证,clear设置为true时，在一次验证后将会删除记录
	if !store.Verify(passwordLoginForms.CaptchaId, passwordLoginForms.Captcha, true) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}

	/*
		登录逻辑
		1.验证手机号是否存在
	*/
	if rsp, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForms.Mobile,
	}); err != nil {
		fmt.Println(err.Error())
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusNotFound, map[string]string{
					"mobile": "手机号未注册",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登入失败,手机号格式错误",
				})
			}
			return
		}
	} else {
		//只是查询到了用户，并没有检查到密码
		if passRsp, passErr := global.UserSrvClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			PassWord:          passwordLoginForms.Password,
			EncryptedPassword: rsp.PassWord,
		}); passErr != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]string{
				"password": "登入失败",
			})
		} else {
			if passRsp.Success {
				//生成 Token
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					uint(rsp.Id),
					rsp.NickName,
					uint(rsp.Role),
					jwt.RegisteredClaims{
						NotBefore: jwt.NewNumericDate(time.Now()),                         //签名的生效时间
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), //token有效期为一周
						Issuer:    "myself",                                               //token发行人
					},
				}
				token, err := j.CreateToken(claims)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成token失败",
					})
					return
				}

				ctx.JSON(http.StatusOK, gin.H{
					"id":         rsp.Id,
					"nick_name":  rsp.NickName,
					"token":      token,
					"expired_at": time.Now().Add(time.Hour * 24 * 7).Unix(),
				})
			} else {
				ctx.JSON(http.StatusBadRequest, map[string]string{
					"msg": "登入失败，密码错误",
				})
			}
		}
	}
}

/*
用户注册
*/
func Register(ctx *gin.Context) {
	//验证注册信息表单格式是否正确，使用了第三方包 validator ,内置验证规则(正则)
	registerForm := forms.RegisterForm{}
	if err := ctx.ShouldBindJSON(&registerForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	fmt.Println(global.ServerConfig.RedisInfo.Host)
	//验证码验证，创建一个redis远程客户端
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	//获取验证信息
	v, err := rdb.Get(context.Background(), registerForm.Mobile).Result()
	if err == redis.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "验证码错误，验证码不存在",
		})
		return
	} else {
		if v != registerForm.Code {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code": "验证码错误",
			})
			return
		}

		//创建用户
		user, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
			NickName: registerForm.Mobile,
			PassWord: registerForm.Password,
			Mobile:   registerForm.Mobile,
		})

		if err != nil {
			zap.S().Errorf("[Register] 注册 【创建用户失败】:%s", err.Error())
			HandleGrpcErrorToHttp(err, ctx)
		}

		j := middlewares.NewJWT()
		claims := models.CustomClaims{
			uint(user.Id),
			user.NickName,
			uint(user.Role), //设置管理权限
			jwt.RegisteredClaims{
				NotBefore: jwt.NewNumericDate(time.Now()),                         //签名的生效时间
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), //token有效期为一周
				Issuer:    "myself",                                               //token发行人
			},
		}
		token, err := j.CreateToken(claims)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": "生成token失败",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"id":         user.Id,
			"nick_name":  user.NickName,
			"token":      token,
			"expired_at": time.Now().Add(time.Hour * 24 * 7).Unix(),
		})
	}
}
