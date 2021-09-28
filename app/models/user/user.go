package user

import (
	"go_blog/app/models"
	"go_blog/pkg/logger"
	"go_blog/pkg/model"
)

type User struct {
	models.BaseModel
	Name     string `gorm:"column:name;type:varchar(255);not null;unique"`
	Email    string `gorm:"column:email;type:varchar(255);default:NULL;unique;"`
	Password string `gorm:"column:password;type:varchar(255)"`
}

func (user *User) Create() error {
	result := model.DB.Create(&user)
	if err := result.Error; err != nil {
		logger.LogError(err)
		return err
	}
	return nil
}
