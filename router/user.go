package router

import (
	"GopherAI/controller/user"

	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(r *gin.RouterGroup) {
	{
		r.POST("/register", user.Register) //post是http协议中的一种请求方法，表示向服务器提交数据，通常用于创建资源或者进行登录等操作
		r.POST("/login", user.Login)
		r.POST("/captcha", user.HandleCaptcha)
	}
}
