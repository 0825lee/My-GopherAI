package jwt

import (
	"GopherAI/common/code"
	"GopherAI/controller"
	"GopherAI/utils/myjwt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 读取jwt
func Auth() gin.HandlerFunc { // 返回值是一个gin.HandlerFunc，这个函数就是定义的中间件函数
	return func(c *gin.Context) {
		res := new(controller.Response)
		// 从请求头中获取 token，支持 Bearer 方式和 URL 参数方式
		var token string
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// 兼容 URL 参数传 token
			token = c.Query("token")
		}

		if token == "" { // 没有找到 token，返回错误
			c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidToken))
			c.Abort()
			return
		}

		log.Println("token is ", token)
		userName, ok := myjwt.ParseToken(token) // 解析 token，获取用户名
		if !ok {                                // 解析失败，返回错误
			c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidToken))
			c.Abort()
			return
		}
		// 解析成功，将用户名存到上下文中，供后续处理函数使用
		c.Set("userName", userName)
		c.Next()
	}
}
