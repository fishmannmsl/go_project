package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_project/fish_farm/fish-api/user-web/api"
	"go_project/fish_farm/fish-api/user-web/middlewares"
)

/*
由于登陆，注册页面不需要已登录，所以不采用一组为单位配置拦截器，而是单个路由来配置
UserRouter := Router.Group("user").Use(middlewares.JWTAuth())
*/
func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	zap.S().Info("配置用户相关的url")
	{
		UserRouter.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		UserRouter.POST("pwd_login", api.PassWordLogin)
		UserRouter.POST("register", api.Register)
	}
}
