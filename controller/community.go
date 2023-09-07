package controller

import (
	"example.com/m/v2/logic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// -----社区相关
func Community(c *gin.Context) {
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
