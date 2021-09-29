package controllers

import (
	"go_blog/pkg/view"
	"net/http"
)

type LoginController struct {
}

func (*LoginController) Login(w http.ResponseWriter, r *http.Request) {
	view.AuthRender(w, view.D{}, "auth.login")
}

func (*LoginController) doLogin(w http.ResponseWriter, r *http.Request) {
	view.AuthRender(w, view.D{}, "auth.login")
}