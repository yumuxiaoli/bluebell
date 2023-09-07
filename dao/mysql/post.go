package mysql

import "example.com/m/v2/models"

func CreatePost(p *models.Post) (err error) {
	err = DB.Create(&p).Error
	return
}
