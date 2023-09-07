package mysql

import "example.com/m/v2/models"

func CreatePost(p *models.Post) (err error) {
	err = DB.Create(&p).Error
	return
}

func GetPostById(pid int64) (post *models.Post, err error) {
	err = DB.Find(&post).Where("post_id = ?", pid).Error
	return
}
