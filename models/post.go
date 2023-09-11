package models

import "time"

type Post struct {
	PostID      int64     `json:"post_id"`
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

// ApiPostDetail 帖子详情接口的结构体
type ApiPostDetail struct {
	AuthorName string             `josn:"author_name"`
	VoteNum    int64              `json:"vote_num"`
	*Post                         // 嵌入帖子结构体
	*Community `json:"community"` // 嵌入社区信息
}
