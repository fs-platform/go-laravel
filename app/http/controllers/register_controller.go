package controllers

import (
	"fmt"
	"go_blog/app/models/user"
	"go_blog/app/requests"
	"go_blog/pkg/logger"
	"go_blog/pkg/view"
	"net/http"
)

type RegisterController struct {
}

func (*RegisterController) Register(w http.ResponseWriter, r *http.Request) {
	view.AuthRender(w, view.D{}, "auth.register")
}

func (*RegisterController) DoRegister(w http.ResponseWriter, r *http.Request) {
	name := r.PostFormValue("name")
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	password_confirmation := r.PostFormValue("password_confirmation")
	_user := user.User{
		Name:            name,
		Email:           email,
		Password:        password,
		PasswordConfirm: password_confirmation,
	}
	errs := requests.ValidateRegistrationForm(_user)
	if len(errs) > 0 {
		// 3. 有错误发生，打印数据
		view.AuthRender(w, view.D{
			"Errors": errs,
			"User":   _user,
		}, "auth.register")
	} else {
		err := _user.Create()
		logger.LogError(err)
		if _user.ID > 0 {
			fmt.Fprintf(w, "插入成功，ID 为%d", _user.ID)
		} else {
			fmt.Fprint(w, "创建用户失败，请联系管理员")
		}
	}
}
