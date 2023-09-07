package mysql

import (
	"database/sql"
	"errors"

	"example.com/m/v2/models"
	"go.uber.org/zap"
)

var ErrorCommunityExis = errors.New("未找到")

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
