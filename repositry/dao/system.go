package dao

import (
	"context"
	"easy-drive/repositry/model"
	"gorm.io/gorm"
)

type SystemDao struct {
	*gorm.DB
}

func NewSystemDao(ctx context.Context) *SystemDao {
	return &SystemDao{GetMysqlClient(ctx)}
}

func NewSystemDaoByDB(db *gorm.DB) *SystemDao {
	return &SystemDao{db}
}

func (sd *SystemDao) GetSystemSettingById(id int) (system *model.System, err error) {
	err = sd.DB.Model(&model.System{}).Where("id = ?", id).First(&system).Error
	return
}
