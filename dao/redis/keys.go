package redis

// redis key

// reids key注意使用命名空间的方式，方便查询和拆分

const (
	Prefix             = "bluebell:"
	KeyPostTime        = "post:time"   // zset;帖子及发帖时间
	KeyPostScore       = "post:score"  // zset;帖子及投票的分数
	KeyPostVotedPrefix = "post:woted:" // zset;记录用户及投票的类型
)

// 给redis key加上前缀
func getRedisKey(key string) string {
	return Prefix + key
}
