package logic

import (
	"example.com/m/v2/dao/mysql"
	"example.com/m/v2/models"
)

func GetCommunityList() ([]*models.Community, error) {
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (*models.Community, error) {
	return mysql.GetCommunityDetailByID(id)
}
