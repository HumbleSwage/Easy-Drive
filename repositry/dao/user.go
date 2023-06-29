package dao

import (
	"context"
	"easy-drive/repositry/model"
	"errors"
	"gorm.io/gorm"
)

type UserDao struct {
	*gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{GetMysqlClient(ctx)}
}

func NewUserDaoByDB(db *gorm.DB) *UserDao {
	return &UserDao{db}
}

func (ud *UserDao) IsEmailExists(email string) (bool, error) {
	var user model.User
	err := ud.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 记录不存在
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}

func (ud *UserDao) IsUserNameExists(userName string) bool {
	var user model.User
	result := ud.DB.Model(&model.User{}).Where("user_name = ?", userName).First(&user)
	err := result.Error
	if err == gorm.ErrRecordNotFound {
		return false
	}
	return true
}

func (ud *UserDao) AddUser(u *model.User) (err error) {
	err = ud.DB.Model(&model.User{}).Create(&u).Error
	return err
}

func (ud *UserDao) GetUserByEmail(email string) (user *model.User, err error) {
	err = ud.DB.Model(&model.User{}).Where("email = ?", email).First(&user).Error
	return
}

func (ud *UserDao) UpdateUser(user *model.User) (err error) {
	err = ud.DB.Model(&model.User{}).Where("user_id = ?", user.UserId).Save(user).Error
	return
}

func (ud *UserDao) GetUserByUserId(userId string) (user *model.User, err error) {
	err = ud.DB.Model(&model.User{}).Where("user_id = ?", userId).First(&user).Error
	return
}
