package middlewares

import (
	"github.com/gin-gonic/gin"
	"go_project/fish_farm/fish-api/user-web/models"
	"net/http"
)

/*
验证用户是否为管理员
*/
func IsAdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		currentUser := claims.(*models.CustomClaims)

		if currentUser.AuthorityId != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg": "无权限",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
