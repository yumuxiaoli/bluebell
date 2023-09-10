package mysql

import (
	"strings"

	"example.com/m/v2/models"
)

func CreatePost(p *models.Post) (err error) {
	err = DB.Create(&p).Error
	return
}

// GetPostById 根据id查询单个帖子数据
func GetPostById(pid int64) (post *models.Post, err error) {
	err = DB.Find(&post).Where("post_id = ?", pid).Error
	return
}

// GetPostList 查询帖子列表函数
func GetPostList(page, size int64) (post []*models.Post, err error) {
	err = DB.Offset(int(page)).Limit(int(size)).Order("create_time desc").Find(&post).Error
	return
}

// 根据给定的id列表查询帖子数据
func GetPostByIDs(ids []string) (postList []*models.Post, err error) {
	err = DB.Where("post_id in (?)", ids).Order("FIND_IN_SET(post_id,'" + strings.Join(ids, ",") + "')").Find(&postList).Error
	return
}
