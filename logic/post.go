package logic

import (
	"time"

	"example.com/m/v2/dao/mysql"
	"example.com/m/v2/models"
	"example.com/m/v2/pkg/snowflake"
)

func CreatePost(p *models.Post) (err error) {
	// 生成post id
	p.ID = snowflake.GenID()
	p.CreateTime = time.Now()
	// 保存到数据库
	return mysql.CreatePost(p)
	// 返回
}
