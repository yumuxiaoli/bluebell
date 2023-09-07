package logic

import (
	"time"

	"example.com/m/v2/dao/mysql"
	"example.com/m/v2/models"
	"example.com/m/v2/pkg/jwt"
	"example.com/m/v2/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) (err error) {
	// 判断用户是否存在
	err = mysql.CheckUserExist(p.Username)
	if err != nil {
		return err
	}
	// 生成UUID
	userID := snowflake.GenID()
	// 构建一个User实例
	u := &models.User{
		UserId:     userID,
		Username:   p.Username,
		Password:   p.Password,
		Email:      p.Email,
		Gender:     p.Gender,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	// 保存数据库
	return mysql.InsertUser(u)
}

func Login(p *models.ParamLogin) (token string, err error) {
	user := &models.ParamLogin{
		Username: p.Username,
		Password: p.Password,
	}
	// 传递的是指针，就能拿到user.UserID
	u, err := mysql.Login(user)
	if err != nil {
		return "", err
	}
	// 生成JWT
	return jwt.GenToken(u.UserId, p.Username)
}
