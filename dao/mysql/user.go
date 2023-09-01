package mysql

import (
	"crypto/md5"
	"encoding/hex"
	"errors"

	"example.com/m/v2/models"
)

var secret string = "2816083598"

// 根据username查询数据
func CheckUserExist(username string) (err error) {
	var user *models.User
	err = DB.Where("username = ?", username).Find(&user).Error
	if err != nil {
		return err
	}
	if user.UserId != 0 {
		return errors.New("用户已存在")
	}
	return nil
}

// 将user存入数据库
func InsertUser(user *models.User) (err error) {
	// 对密码进行加密
	user.Password = encryptPassword(user.Password)
	err = DB.Create(&user).Error
	return err
}

func encryptPassword(password string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(password)))
}
