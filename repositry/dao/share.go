package dao

import (
	"context"
	"easy-drive/repositry/model"
	"gorm.io/gorm"
)

type ShareDao struct {
	*gorm.DB
}

func NewShareDao(ctx context.Context) *ShareDao {
	return &ShareDao{GetMysqlClient(ctx)}
}

func NewShareDaoByDB(db *gorm.DB) *ShareDao {
	return &ShareDao{db}
}

func (sd *ShareDao) GetShareInfoById(shareId string) (share *model.Share, err error) {
	err = sd.DB.Model(&model.Share{}).Where("share_id = ?", shareId).First(&share).Error
	return
}

func (sd *ShareDao) GetShareFileByUserId(userId string, pageNum, pageSize int) (shares []*model.Share, count int64, err error) {
	err = sd.DB.Model(&model.Share{}).Where("user_id = ?", userId).
		Offset((pageNum - 1) * pageSize).Limit(pageSize).
		Find(&shares).Order("share_time DESC").Error
	if err != nil {
		return nil, 0, err
	}

	err = sd.DB.Model(&model.Share{}).Where("user_id = ?", userId).Count(&count).Error
	return
}

func (sd *ShareDao) DeleteShareFile(shareId, userId string) (err error) {
	err = sd.DB.Model(&model.Share{}).Where("share_id = ? AND user_id = ?", shareId, userId).Delete(&model.Share{}).Error
	return
}
