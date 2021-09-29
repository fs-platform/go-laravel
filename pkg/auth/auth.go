package auth

import (
	"errors"
	"go_blog/app/models/user"
	"go_blog/pkg/session"
	"gorm.io/gorm"
)

func _getUID() uint64 {
	_uid := session.Get("uid")
	if _uid == nil {
		return 0
	}
	uid := _uid.(uint64)
	if uid > 0 {
		return uid
	}
	return 0
}

func User() user.User {
	uid := _getUID()
	if uid > 0 {
		_user, err := user.Get(uid)
		if err == nil {
			return _user
		}
	}
	return user.User{}
}

func Attempt(email string, password string) error {
	// 1. 根据 Email 获取用户
	_user, err := user.GetByEmail(email)
	// 2. 如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("账号不存在或密码错误")
		} else {
			return errors.New("内部错误，请稍后尝试")
		}
	}
	// 3. 匹配密码
	if !_user.ComparePassword(password) {
		return errors.New("密码不正确")
	}

	// 4. 登录用户，保存会话
	session.Put("uid", _user.ID)

	return nil
}

// Login 登录指定用户
func Login(_user user.User) {
	session.Put("uid", _user.ID)
}

// Logout 退出用户
func Logout() {
	session.Forget("uid")
}

// Check 检测是否登录
func Check() bool {
	return _getUID() > 0
}
