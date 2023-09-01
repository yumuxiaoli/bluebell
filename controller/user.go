package controller

import (
	"fmt"
	"net/http"

	"example.com/m/v2/logic"
	"example.com/m/v2/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 处理注册请求的函数
func Register(c *gin.Context) {
	// 1、参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回相应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"msg":  "参数错误",
		})
		return
	}
	// 对请求参数进行详细业务规则校验
	if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.RePassword != p.Password {
		zap.L().Error("SignUp with invalid param")
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"msg":  "参数错误",
		})
		return
	}
	fmt.Println(p)
	// 2、业务处理
	logic.SignUp(p)
	// 3、返回响应
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
}
