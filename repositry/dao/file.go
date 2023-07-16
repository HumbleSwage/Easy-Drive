package dao

import (
	"context"
	"easy-drive/repositry/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
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

func (fd *FileDao) SelectFile(pid string, pageNum, pageSize int) (files []*model.File, count int64, err error) {
	err = fd.DB.Model(&model.File{}).Where("flag = 2 AND parent_id = ?", pid).Find(&files).
		Offset((pageNum - 1) * pageSize).Limit(pageSize).Error
	if err != nil {
		return nil, 0, err
	}
	err = fd.DB.Model(&model.File{}).Where("flag = 2 AND parent_id = ?", pid).Count(&count).Error
	return
}

func (fd *FileDao) SelectFileByFuzzyName(pid, fileName string, pageNum, pageSize int) (files []*model.File, count int64, err error) {
	err = fd.DB.Model(&model.File{}).Where("parent_id = ?", pid).
		Where("name LIKE ?", strings.Join([]string{"%", fileName, "%"}, "")).Find(&files).Error
	if err != nil {
		return nil, 0, err
	}
	err = fd.DB.Model(&model.File{}).
		Where("name LIKE ?", strings.Join([]string{"%", fileName, "%"}, "")).
		Count(&count).Error
	return
}

func (fd *FileDao) GetUserFilesByCategory(user *model.User, filePid string, category int, pageNum, pageSize int) (files []*model.File, count int64, err error) {
	err = fd.DB.Model(&model.File{}).Order("updated_at DESC").Where("user_id = ?", user.UserId).Where("parent_id = ?", filePid).Where("category = ? ", category).Where("is_directory = 0 AND flag = 2").
		Offset((pageNum - 1) * pageSize).Limit(pageSize).
		Find(&files).Error
	if err != nil {
		return nil, 0, err
	}
	err = fd.DB.Model(&model.File{}).Order("updated_at DESC").Where("user_id = ?", user.UserId).Where("parent_id = ?", filePid).Where("category = ? ", category).Where("is_directory = 0 AND flag = 2").
		Count(&count).Error
	return
}

func (fd *FileDao) GetFileById(fileId, userId string) (file *model.File, err error) {
	err = fd.DB.Model(&model.File{}).Where("file_id = ? AND user_id = ?", fileId, userId).First(&file).Error
	return
}

func (fd *FileDao) GetFileByFileId(fileId string) (file *model.File, err error) {
	err = fd.DB.Model(&model.File{}).Where("file_id = ?", fileId).Find(&file).Error
	return
}

func (fd *FileDao) GetFileByFileName(fileName, filePid, userId string) (file *model.File, err error) {
	err = fd.DB.Model(&model.File{}).Where("name = ? AND parent_id = ? AND flag = 2 AND user_id = ?", fileName, filePid, userId).First(&file).Error
	return
}

func (fd *FileDao) GetUserFiles(userId, filePid string, pageNum, pageSize int) (files []*model.File, count int64, err error) {
	err = fd.DB.Model(&model.File{}).Order("updated_at DESC").Where("user_id = ? AND parent_id = ? AND flag = 2", userId, filePid).
		Offset((pageNum - 1) * pageSize).Limit(pageSize).
		Find(&files).Error
	if err != nil {
		return nil, 0, err
	}
	err = fd.DB.Model(&model.File{}).Where("user_id = ? AND parent_id = ?  AND flag = 2", userId, filePid).
		Count(&count).Error
	return
}

func (fd *FileDao) GetFiles(userId string, fileIds []string) (files []*model.File, err error) {
	err = fd.DB.Model(&model.File{}).Where("user_id = ?", userId).Where("file_id IN ?", fileIds).Find(&files).Error
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

func (fd *FileDao) GetFileInSameFolder(userId string, folderPath []string) (files []*model.File, err error) {
	err = fd.DB.Model(&model.File{}).Where("user_id = ? AND is_directory = 0", userId).Where("parent_id IN ?", folderPath).Clauses(
		clause.OrderBy{
			Expression: clause.Expr{SQL: "FIELD(file_id,?)", Vars: []interface{}{folderPath}, WithoutParentheses: true},
		},
	).Find(&files).Error
	return
}

func (fd *FileDao) GetFolderInfo(userId string, folderPath []string) (folders []*model.File, err error) {
	err = fd.DB.Model(&model.File{}).Where("user_id = ? AND is_directory = 1", userId).Where("file_id IN ?", folderPath).Clauses(
		clause.OrderBy{
			Expression: clause.Expr{SQL: "FIELD(file_id,?)", Vars: []interface{}{folderPath}, WithoutParentheses: true},
		},
	).Find(&folders).Error
	return
}

func (fd *FileDao) GetFileInFolder(pid string) (files []*model.File, err error) {
	err = fd.DB.Model(&model.File{}).Where("parent_id = ?", pid).Find(&files).Error
	return
}

func (fd *FileDao) SelectFolderInfo(folderPath []string) (folders []model.File, err error) {
	err = fd.DB.Model(&model.File{}).Where("file_id IN ?", folderPath).Where("is_directory = 1").Clauses(
		clause.OrderBy{
			Expression: clause.Expr{SQL: "FIELD(file_id,?)", Vars: []interface{}{folderPath}, WithoutParentheses: true},
		},
	).Find(&folders).Error
	return
}

func (fd *FileDao) CheckFolderName(userId, filePid, fileName string) (count int64, err error) {
	err = fd.DB.Model(&model.File{}).Where("user_id = ? AND parent_id = ? AND name = ? AND is_directory = 1", userId, filePid, fileName).Count(&count).Error
	return

}

func (fd *FileDao) GetFolder(userId, filePid string, fileId ...[]string) (files []*model.File, err error) {
	if fileId != nil {
		err = fd.DB.Model(&model.File{}).Where("user_id = ? AND flag = 2", userId).Where("is_directory = 1 AND parent_id = ?", filePid).Where("file_id NOT IN ?", fileId).Find(&files).Error
	} else {
		err = fd.DB.Model(&model.File{}).Where("user_id = ? AND flag = 2", userId).Where("is_directory = 1 AND parent_id = ?", filePid).Find(&files).Error
	}
	return
}

func (fd *FileDao) GetDelFile(userId string, pageNum, pageSize int) (files []*model.File, count int64, err error) {
	err = fd.DB.Model(&model.File{}).Where("user_id = ? AND flag = 1", userId).
		Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&files).Error
	if err != nil {
		return nil, 0, err
	}
	err = fd.DB.Model(&model.File{}).Where("user_id = ? AND flag = 1", userId).
		Offset((pageNum - 1) * pageSize).Limit(pageSize).Count(&count).Error
	return
}
