package models

import "time"

type Community struct {
	Id            uint      `json:"id"`
	CommunityId   int64     `json:"community_id"`
	CommunityName string    `json:"community_name"`
	Introduction  string    `json:"introduction"`
	CreateTime    time.Time `json:"create_time"`
	UpdateTime    time.Time `json:"update_time"`
}

func (Community) TableName() string {
	return "community"
}

type ResCommunity struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
