package controller

import (
	"strconv"

	"example.com/m/v2/logic"
	"example.com/m/v2/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreatePost(c *gin.Context) {
	// 获取参数及参数的校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON(p) error", zap.Any("err", err))
		zap.L().Error("create post with invalid param")
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 从c取到当前发送的用户的id
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	// 2.创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3.返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetail 获取帖子详情的处理函数
func GetPostDetail(c *gin.Context) {
	// 获取参数
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 根据id去除帖子数据(查数据库)
	data, err := logic.GetPostById(pid)
	if err != nil {
		zap.L().Error("logic.GetPostById(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// 获得帖子列表的处理函数
func GetPostList(c *gin.Context) {
	// 获取分页参数
	offsetStr := c.Query("page")
	limitStr := c.Query("size")

	page, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		page = 0
	}
	size, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		size = 10
	}
	// 获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// 根据前端传来的参数动态的获取帖子列表
// 安装创建时间排序，或者 按照 分数排序
// 根据id去数据库查询帖子详细信息
func GetPostList2(c *gin.Context) {
	// 获取分页参数

	// 初始化结构体参数时指定初始承诺书
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}

	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostList2 with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取数据
	data, err := logic.GetPostList2(p)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 返回响应
	ResponseSuccess(c, data)
}
