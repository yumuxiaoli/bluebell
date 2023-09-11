package redis

import (
	"fmt"
	"strconv"
	"time"

	"example.com/m/v2/models"
	"github.com/go-redis/redis"
)

func getIDsFormKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * page
	end := start + size - 1

	// ZREVRANGE 按分数从大到小的顺序查询指定数量的元素
	return rdb.ZRevRange(key, start, end).Result()
}

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 从redis获取id
	key := getRedisKey(KeyPostTime)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScore)
	}
	fmt.Println("======================\n", p.Order, "==", key)
	return getIDsFormKey(key, p.Page, p.Size)
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

// GetCommunityPostIDsInOrder 根据社区查询ids
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	orderKey := getRedisKey(KeyPostTime)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScore)
	}
	fmt.Println("======================\n", p.Order, "==", orderKey)
	// 使用zinterstore 把分区的帖子set余帖子分数的zset 生成一个新zset
	// 针对新的zset 按之前的梁洛级取数据

	// 社区的key
	cKey := getRedisKey(KeyCommunityPrefix + strconv.Itoa(int(p.CommunityID)))
	// 利用缓存Key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if rdb.Exists(orderKey).Val() < 1 {
		// 不存在，需要计算
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, cKey, orderKey) // zinterstore 计算
		pipeline.Expire(key, 60*time.Second) // 设置超时时间
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	return getIDsFormKey(key, p.Page, p.Size)
}
