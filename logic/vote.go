package logic

import (
	"strconv"

	"example.com/m/v2/dao/redis"
	"example.com/m/v2/models"
	"go.uber.org/zap"
)

// 本项目使用简化版的投票分数

// VoteForPost 为帖子投票的函数
func VoteForPost(userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost",
		zap.Int64("userID", userID),
		zap.String("postID", p.PostID),
		zap.Int8("Direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}
