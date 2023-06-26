package router

import (
	"github.com/gin-gonic/gin"
	
	"go_project/fish_farm/fish-api/oss-web/handler"
)

func InitOssRouter(Router *gin.RouterGroup) {
	OssRouter := Router.Group("oss")
	{
		//OssRouter.GET("token", middlewares.JWTAuth(), middlewares.IsAdminAuth(), handler.Token)
		OssRouter.GET("token", handler.Token)
		OssRouter.POST("/callback", handler.HandlerRequest)
	}
}
