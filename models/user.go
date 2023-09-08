package models

import "time"

type User struct {
	Id         uint      `json:"id"`
	UserId     int64     `json:"user_id"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Email      string    `json:"email"`
	Gender     bool      `json:"gender"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

func (User) TableName() string {
	return "user"
}

type ResData struct {
	UserID      int64  `json:"userID,string"`
	Username    string `json:"username"`
	AccessToken string `json:"accessToken"`
}
