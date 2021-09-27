package controllers

import (
	"fmt"
	"net/http"
)

type ArticlesController struct {
}

func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {

}

func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {

}

func (*ArticlesController) Update(w http.ResponseWriter, r *http.Request) {

}

func (*ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {

}

func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {

}

func (*ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {

}

func (*ArticlesController) About(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:Aron@example.com\">Aron@example.com</a>")
}

// NotFound 404 页面
func (*ArticlesController) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}
