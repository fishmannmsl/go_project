package initialize

/*
作为一个初始化所有服务函数的包
减少main.go内的代码量
*/

import (
	"github.com/gin-gonic/gin"
	"go_project/fish_farm/fish-api/goods-web/middlewares"
	"go_project/fish_farm/fish-api/goods-web/router"
)

/*
初始化路由
*/
func Routers() *gin.Engine {
	Router := gin.Default()
	//配置跨域处理
	Router.Use(middlewares.Cors())
	ApiGroup := Router.Group("/u/v1")
	//初始化路由组
	router.InitUserRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup)
	return Router
}
