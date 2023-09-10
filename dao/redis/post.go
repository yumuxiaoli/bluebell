package redis

import "example.com/m/v2/models"

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
