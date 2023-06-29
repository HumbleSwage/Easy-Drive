package dao

import (
	"context"
	"easy-drive/repositry/model"
	"gorm.io/gorm"
)

type FileDao struct {
	*gorm.DB
}

func NewFileDao(ctx context.Context) *FileDao {
	return &FileDao{GetMysqlClient(ctx)}
}

func NewFileDaoByDB(db *gorm.DB) *FileDao {
	return &FileDao{db}
}

func (fd *FileDao) CreateFile(file *model.File) (err error) {
	err = fd.DB.Model(&model.File{}).Create(&file).Error
	return
}

func (fd *FileDao) GetUserFilesByCategory(user *model.User, category, pageNum, pageSize int, count int64) (files []*model.File, err error) {
	err = fd.DB.Model(&model.File{}).Order("updated_at DESC").Where("user_id = ? AND category = ??", user.UserId, category).
		Offset((pageNum - 1) * pageSize).Limit(pageSize).
		Find(&files).Error
	if err != nil {
		return nil, err
	}
	err = fd.DB.Model(&model.File{}).Where("user_id = ? AND category = ? AND status = 0", user.UserId, category).
		Offset((pageNum - 1) * pageSize).Limit(pageSize).
		Count(&count).Error
	return
}

func (fd *FileDao) GetUserFiles(user *model.User, pageNum, pageSize int, count int64) (files []*model.File, err error) {
	err = fd.DB.Model(&model.File{}).Order("updated_at DESC").Where("user_id = ?", user.UserId).
		Offset((pageNum - 1) * pageSize).Limit(pageSize).
		Find(&files).Error
	if err != nil {
		return nil, err
	}
	err = fd.DB.Model(&model.File{}).Where("user_id = ?", user.UserId).
		Offset((pageNum - 1) * pageSize).Limit(pageSize).
		Count(&count).Error
	return
}

func (fd *FileDao) SelectUserFileSpace(user *model.User, status int) (sum int64, err error) {
	err = fd.DB.Model(&model.File{}).Where("user_id = ? AND status = ?", user.UserId, status).Select("IFNULL(SUM(size),0)").Scan(&sum).Error
	return
}

func (fd *FileDao) SelectFileByMd5(md5 string, flag int) (file *model.File, err error) {
	err = fd.DB.Model(&model.File{}).Where("md5 = ? AND flag = ?", md5, flag).First(&file).Error
	return
}

func (fd *FileDao) SelectFileInFolder(pid string, user *model.User, fileName string) (count int64, err error) {
	err = fd.DB.Model(&model.File{}).Where("user_id = ? AND parent_id = ? AND name = ?", user.UserId, pid, fileName).Count(&count).Error
	return
}

func (fd *FileDao) UpdateFile(file *model.File) (err error) {
	err = fd.DB.Model(&model.File{}).Where("file_id = ?", file.FileId).Save(&file).Error
	return
}
