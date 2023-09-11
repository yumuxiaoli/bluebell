package logic

import (
	"time"

	"example.com/m/v2/dao/mysql"
	"example.com/m/v2/dao/redis"
	"example.com/m/v2/models"
	"example.com/m/v2/pkg/snowflake"
	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	// 生成post id
	p.PostID = snowflake.GenID()
	p.CreateTime = time.Now()
	// 保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return
	}
	return redis.CreatePost(p.PostID, p.CommunityID)
	// 返回
}

// GetPostById 根据帖子id查询帖子详情数据
func GetPostById(pid int64) (data *models.ApiPostDetail, err error) {
	post, err := mysql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById(pid) failed", zap.Int64("pid", pid), zap.Error(err))
		return
	}
	// 根据作者id查询作者信息
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		return
	}
	commuity, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
		return
	}
	data = &models.ApiPostDetail{
		AuthorName: user.Username,
		Post:       post,
		Community:  commuity,
	}
	return
}

// // GetPostList 获取帖子列表
// func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
// 	posts, err := mysql.GetPostList(page, size)
// 	if err != nil {
// 		return nil, err
// 	}
// 	data = make([]*models.ApiPostDetail, 0, len(posts))
// 	for _, post := range posts {
// 		// 根据作者id查询作者信息
// 		user, err := mysql.GetUserById(post.AuthorID)
// 		if err != nil {
// 			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
// 				zap.Int64("author_id", post.AuthorID),
// 				zap.Error(err))
// 			return nil, err
// 		}
// 		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
// 		if err != nil {
// 			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
// 				zap.Int64("community_id", post.CommunityID),
// 				zap.Error(err))
// 			return nil, err
// 		}
// 		postdetail := &models.ApiPostDetail{
// 			AuthorName: user.Username,
// 			Post:       post,
// 			Community:  community,
// 		}
// 		data = append(data, postdetail)
// 	}
// 	return
// }

// GetPostList2()
func GetPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 去redis查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	// 根据id去数据库查询帖子详细信息
	posts, err := mysql.GetPostByIDs(ids)
	if err != nil {
		return
	}
	// 提前查询好每篇投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}
	// 将帖子的作者及分区新查询出来填充到帖子中
	for idx, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			return nil, err
		}
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			return nil, err
		}
		postdetail := &models.ApiPostDetail{
			AuthorName: user.Username,
			VoteNum:    voteData[idx],
			Post:       post,
			Community:  community,
		}
		data = append(data, postdetail)
	}
	return
}

func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 去redis查询id列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetCommunityPostIDsInOrder(p) return 0 data")
		return
	}
	// 根据id去数据库查询帖子详细信息
	posts, err := mysql.GetPostByIDs(ids)
	if err != nil {
		return
	}
	// 提前查询好每篇投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}
	// 将帖子的作者及分区新查询出来填充到帖子中
	for idx, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			return nil, err
		}
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			return nil, err
		}
		postdetail := &models.ApiPostDetail{
			AuthorName: user.Username,
			VoteNum:    voteData[idx],
			Post:       post,
			Community:  community,
		}
		data = append(data, postdetail)
	}
	return
}

// GetPostListNew 将连个查询接口合二为一
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 根据请求参数的不同，执行不同的逻辑
	if p.CommunityID == 0 {
		// 查所有
		data, err = GetPostList(p)
	} else {
		// 根据社区id查询
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return
}
