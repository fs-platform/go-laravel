package controllers

import (
	"fmt"
	"go_blog/app/models/article"
	"go_blog/pkg/logger"
	"go_blog/pkg/route"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
)

type ArticlesController struct {
}

func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
	var (
		tmp *template.Template
		err error
	)
	tmp, err = template.ParseFiles("resources/views/articles/index.gohtml")
	logger.LogError(err)
	articles, err := article.GetAll()
	fmt.Println(articles)
	if err != nil {
		// 数据库错误
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 服务器内部错误")
	} else {
		err = tmp.Execute(w, articles)
	}
}

func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	data, err := article.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
			return
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
			return
		}
	}
	tmpl, err := template.ParseFiles("resources/views/articles/show.gohtml")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, data)
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
