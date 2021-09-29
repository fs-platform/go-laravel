package user

import (
	"go_blog/pkg/model"
)

// Get 通过 ID 获取用户
func Get(idstr uint64) (User, error) {
	var user User
	if err := model.DB.First(&user, idstr).Error; err != nil {
		return user, err
	}
	return user, nil
}

// GetByEmail 通过 Email 来获取用户
func GetByEmail(email string) (User, error) {
	var user User
	if err := model.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}