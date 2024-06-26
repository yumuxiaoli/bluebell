package middleware

import (
	"strings"

	"example.com/m/v2/controller"
	"example.com/m/v2/pkg/jwt"
	"example.com/m/v2/utils"
	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 基于JWT认证的中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式1、放在请求头 2、放在请求体 3、放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controller.ResponseError(c, utils.CodeNeedLogin)
			c.Abort()
			return
		}
		pants := strings.SplitN(authHeader, " ", 2)
		pants[0] = strings.Replace(pants[0], " ", "", -1)
		pants[1] = strings.Replace(pants[1], " ", "", -1)

		if !(len(pants) == 2 && pants[0] == "Bearer") {
			controller.ResponseError(c, utils.CodeInvalidToken)
			c.Abort()
			return
		}

		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(pants[1])
		if err != nil {
			controller.ResponseError(c, utils.CodeInvalidToken)
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set(controller.CtxUserIDKey, mc.UserID)

		c.Next() // 后续的处理函数可以通过c.Get("username")来获取当前请求的用户信息
	}
}
