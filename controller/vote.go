package controller

import (
	"example.com/m/v2/logic"
	"example.com/m/v2/models"
	"example.com/m/v2/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func PostVoteControl(c *gin.Context) {
	// 参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, utils.CodeInvalidParam)
			return
		}
		errDate := removeTopStruct(errs.Translate(trans))
		ResponseErrorWithMsg(c, utils.CodeInvalidParam, errDate)
		return
	}
	// 获取当前请求用户的ID
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, utils.CodeNeedLogin)
	}
	err = logic.VoteForPost(userID, p)
	if err != nil {
		zap.L().Error("logic.VoteForPost Invalid", zap.Error(err))
		ResponseError(c, utils.CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
