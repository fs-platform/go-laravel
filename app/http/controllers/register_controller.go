package controllers

import (
	"fmt"
	"go_blog/pkg/view"
	"net/http"
)

type RegisterController struct {
}

func (*RegisterController) Register(w http.ResponseWriter, r *http.Request) {
	view.AuthRender(w, view.D{}, "auth.register")
}

func (*RegisterController) DoRegister(w http.ResponseWriter, r *http.Request) {
	fmt.Println(4)
}
