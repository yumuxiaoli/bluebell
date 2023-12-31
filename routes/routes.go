package routes

import (
	"net/http"

	"example.com/m/v2/controller"
	"example.com/m/v2/logger"
	"example.com/m/v2/middleware"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"

	_ "example.com/m/v2/docs"
	"github.com/gin-contrib/pprof"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // gin设置为发布模式
	}

	r := gin.New()
	// r.Use(logger.GinLogger(), logger.GinRecovery(true), middleware.RateLimitMiddleware(time.Second*2, 1))
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.LoadHTMLFiles("./templates/index.html")
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	// 注册业务路由
	v1.POST("/register", controller.Register)
	v1.POST("/login", controller.Login)

	v1.GET("/posts", controller.GetPostList)

	v1.Use(middleware.JWTAuthMiddleware()) // 认证中间件
	{
		v1.GET("/community", controller.Community)
		v1.GET("/community/:id", controller.CommunityDetail)

		v1.POST("/post", controller.CreatePost)
		v1.GET("/post/:id", controller.GetPostDetail)
		// v1.GET("/posts/", controller.GetPostList)

		v1.POST("/vote", controller.PostVoteControl)
	}

	r.GET("/ping", func(c *gin.Context) {
		//如果是登录用户,判断请求头中是否有 有效的JWT ？
		c.String(http.StatusOK, "pong")
	})
	pprof.Register(r) // 注册pprof相关路径

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"msg":  "not found",
		})
	})

	return r
}
