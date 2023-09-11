package models

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// 定义请求的参数结构体

// register 请求参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Gender     bool   `json:"gender"`
}

// login 请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamVoteData 投票数据
type ParamVoteData struct {
	//UserID
	PostID    string `json:"post_id,string" binding:"required"`       // 帖子id
	Direction int8   `json:"direction,string" binding:"oneof=1 0 -1"` // 赞成票(1)反对票(-1)
}

// ParamPostList 获取帖子列表query string参数
type ParamPostList struct {
	Page        int64  `json:"page"`
	Size        int64  `json:"size"`
	Order       string `json:"order"`
	CommunityID int64  `json:"community_id" form:"community_id"` // 可以为空
}
