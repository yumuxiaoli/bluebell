package mysql

import (
	"crypto/md5"
	"encoding/hex"
	"errors"

	"example.com/m/v2/models"
)

var secret string = "2816083598"
var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("密码错误")
)

// CheckUserExist 检查用户是否存在
func CheckUserExist(username string) (err error) {
	var user *models.User
	err = DB.Where("username = ?", username).Find(&user).Error
	if err != nil {
		return err
	}
	if user.UserId != 0 {
		return ErrorUserExist
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

// encryptPassword 密码加密
func encryptPassword(password string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(password)))
}

func Login(p *models.ParamLogin) (user models.User, err error) {
	err = DB.Where("username = ?", p.Username).Find(&user).Error
	if err != nil {
		return user, err
	}
	if user.UserId == 0 {
		return user, ErrorUserNotExist
	}
	// 判断密码是否争取
	password := encryptPassword(p.Password)
	if password != user.Password {
		return user, ErrorInvalidPassword
	}
	return user, nil
}

// GetUserById 根据id获取用户信息
func GetUserById(uid int64) (user *models.User, err error) {
	err = DB.Find(&user).Where("user_id = ?", uid).Error
	return
}
