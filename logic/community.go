package logic

import (
	"example.com/m/v2/dao/mysql"
	"example.com/m/v2/models"
)

func GetCommunityList() (data []*models.Community, err error) {
	return mysql.GetCommunityList()
}
