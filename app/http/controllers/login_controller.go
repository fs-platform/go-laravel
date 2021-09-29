package controllers

import (
	"fmt"
	"go_blog/pkg/auth"
	"go_blog/pkg/view"
	"net/http"
)

type LoginController struct {
}

func (*LoginController) Login(w http.ResponseWriter, r *http.Request) {
	view.AuthRender(w, view.D{}, "auth.login")
}

func (*LoginController) DoLogin(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	err := auth.Attempt(email, password)
	if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		fmt.Println(err)
		view.AuthRender(w, view.D{
			"Email":    email,
			"Error": err.Error(),
			"Password": password,
		}, "auth.login")
	}
}
