package controller

import (
	"errors"
	"fmt"

	"example.com/m/v2/dao/mysql"
	"example.com/m/v2/logic"
	"example.com/m/v2/models"
	"example.com/m/v2/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// 处理注册请求的函数
func Register(c *gin.Context) {
	// 参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回相应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, utils.CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, utils.CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 业务处理
	err := logic.SignUp(p)
	if err != nil {
		zap.L().Error("logic.resgister failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, utils.CodeUserExist)
			return
		}
		ResponseError(c, utils.CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, utils.CodeSuccess)
}

func Login(c *gin.Context) {
	// 获取请求参数及参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回相应
		zap.L().Error("Login with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, utils.CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, utils.CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 业务逻辑处理
	data, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed:", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, utils.CodeUserNotExist)
			return
		}
		ResponseError(c, utils.CodeInvalidPassword)
		return
	}
	// 返回响应
	ResponseSuccess(c, gin.H{
		"accessToken": data.AccessToken,
		"username":    data.Username,
		"userID":      fmt.Sprintf("%d", data.UserID),
	})
}
