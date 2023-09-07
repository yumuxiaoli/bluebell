package mysql

import (
	"database/sql"
	"errors"

	"example.com/m/v2/models"
	"go.uber.org/zap"
)

var ErrCommunityExis = errors.New("未找到")
var ErrInvalidID = errors.New("未找到")

func GetCommunityList() (community []*models.Community, err error) {
	err = DB.Find(&community).Error
	if err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
		}
		return nil, err
	}
	return
}

func GetCommunityDetailByID(id int64) (community *models.Community, err error) {
	err = DB.Find(&community).Where("community_id=?", id).Error
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrInvalidID
		}
	}
	return
}
