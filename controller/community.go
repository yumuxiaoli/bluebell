package controller

import (
	"strconv"

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

// CommunityDetail 社区分类详情
func CommunityDetail(c *gin.Context) {
	// 获取社区id
	idStr := c.Param("id")
	// 查询到所有的社区，(community_id,community_name)以列表的形式返回
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
