package redis

import (
	"fmt"

	"example.com/m/v2/settings"
	"github.com/go-redis/redis"
)

// 声明一个全局的rdb变量
var (
	rdb *redis.Client
	Nil = redis.Nil
)

// 初始化连接
func Init(cfg *settings.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,
			cfg.Port,
		),
		Password: cfg.Password,
		DB:       cfg.Db,
		PoolSize: cfg.PoolSize,
	})
	_, err = rdb.Ping().Result()
	if err != nil {
		return err
	}
	return
}

func Close() {
	_ = rdb.Close()
}
