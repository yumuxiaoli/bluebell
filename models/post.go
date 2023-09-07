package models

import "time"

type Post struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title" binding:"required"`
	Content     string    `json:"content" binding:"required"`
	AuthorID    int64     `json:"author_id"`
	CommunityID int64     `json:"community_id" binding:"required"`
	Status      int32     `json:"status"`
	CreateTime  time.Time `json:"create_time"`
}

func (Post) TableName() string {
	return "post"
}
