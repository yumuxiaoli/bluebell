package redis

import (
	"errors"
	"math"
	"time"

	"github.com/go-redis/redis"
)

/*
投票的几种情况：
direction = 1时，有两种情况
	1.之前没有投票，现在投赞成票		--> 跟新分数和投票记录	差值的绝对值:1 +432
	2.之前投反对票，现在改投赞成票		--> 跟新分数和投票记录	差值的绝对值:2 +432*2
direction = 0时，有两种情况
	1.之前投赞成票，现在要取消投票		--> 跟新分数和投票记录	差值的绝对值:1 -432
	2.之前投反对票，现在要取消投票		--> 跟新分数和投票记录	差值的绝对值:1 +432
direction = -1时,有两种情况
	1.之前没有投票，现在投反对票		--> 跟新分数和投票记录 差值的绝对值:1 -432
	2.之前投赞成票，现在投反对票		--> 跟新分数和投票记录 差值的绝对值:2 -432*2

投票的限制：
	每个贴子发表之日起一个星期内润许用户投票，超过一个星期就不允许投票了
	1.到期之后redis中保存的赞成票数以及反对票数存储到mysql表中
	2.到期之后删除那个 KeyPostVotedZSetPF
*/

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每一票的值
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
)

func CreatePost(postID int64) (err error) {
	pipeline := rdb.TxPipeline()
	// 帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTime), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScore), redis.Z{
		Score:  float64(0),
		Member: postID,
	})
	_, err = pipeline.Exec()

	return err
}

func VoteForPost(userID, postID string, value float64) (err error) {
	// 1.判断投票限制
	// 去redis取
	postTime := rdb.ZScore(getRedisKey(KeyPostTime), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	// 2和3需要梵高一个pipeline中
	// 2.更新帖子的分数
	// 先查之前的投票纪录
	ov := rdb.ZScore(getRedisKey(KeyPostVotedPrefix+postID), userID).Val()
	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value) // 计算两次投票的差值
	pipeline := rdb.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScore), op*diff*scorePerVote, postID).Result()
	if err != nil {
		return err
	}
	// 3.记录用户为该帖子投票的数据
	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedPrefix+postID), postID).Result()
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedPrefix+postID), redis.Z{
			Score:  value,
			Member: userID,
		}).Result()
	}
	_, err = pipeline.Exec()
	return
}
