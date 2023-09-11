package redis

import (
	"example.com/m/v2/models"
	"github.com/go-redis/redis"
)

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 从redis获取id
	key := getRedisKey(KeyPostTime)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScore)
	}

	// 确定查询的索引起始点
	start := (p.Page - 1) * p.Page
	end := start + p.Size - 1

	// ZREVRANGE 按分数从大到小的顺序查询指定数量的元素
	return rdb.ZRevRange(key, start, end).Result()
}

// GetPostVoteData 根据ids查询每篇帖子的投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	// for _, id := range ids {
	// 	key := getRedisKey(KeyPostVotedPrefix + id)
	// 	// 查找key中分数是1的元素的数量
	// 	v := rdb.ZCount(key, "1", "1").Val()
	// 	data = append(data, v)
	// }

	// 使用pipeleine一次发送多条命令，减少RTT
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedPrefix + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}
