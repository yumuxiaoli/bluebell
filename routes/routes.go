package routes

import (
	"net/http"

	"example.com/m/v2/controller"
	"example.com/m/v2/logger"
	"example.com/m/v2/middleware"
	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // gin设置为发布模式
	}

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册业务路由
	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r.GET("/ping", middleware.JWTAuthMiddleware(), func(c *gin.Context) {
		//如果是登录用户,判断请求头中是否有 有效的JWT ？
		c.String(http.StatusOK, "pong")
	})
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"msg":  "not found",
		})
	})

	return r
}
