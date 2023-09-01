package logic

import (
	"example.com/m/v2/dao/mysql"
	"example.com/m/v2/models"
	"example.com/m/v2/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) {
	// 判断用户是否存在
	mysql.QueryUserByUserName()
	// 生成UUID
	snowflake.GenID()
	// 保存数据库
	mysql.InsertUser()
}
