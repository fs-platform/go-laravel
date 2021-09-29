package user

import (
	"go_blog/app/models"
	"go_blog/pkg/logger"
	"go_blog/pkg/model"
)

type User struct {
	models.BaseModel
	Name            string `gorm:"type:varchar(255);not null;unique" valid:"name"`
	Email           string `gorm:"type:varchar(255);default:NULL;unique;" valid:"email"`
	Password        string `gorm:"type:varchar(255)" valid:"password"`
	PasswordConfirm string `gorm:"-" valid:"password_confirm"`
}

func (user *User) Create() error {
	result := model.DB.Create(&user)
	if err := result.Error; err != nil {
		logger.LogError(err)
		return err
	}
	return nil
}

// ComparePassword 对比密码是否匹配
func (user User) ComparePassword(password string) bool {
	return user.Password == password
}
