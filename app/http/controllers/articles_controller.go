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
	"unicode/utf8"
)

// ArticlesFormData 创建博文表单数据
type ArticlesFormData struct {
	Title, Body string
	URL         string
	Errors      map[string]string
}

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

func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
	storeUrl := route.RouteName2URL("articles.store")
	tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}
	data := new(ArticlesFormData)
	data.URL = storeUrl
	err = tmpl.Execute(w, data)
	logger.LogError(err)
}

func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		// 解析错误，这里应该有错误处理
		fmt.Fprint(w, "请提供正确的数据！")
		return
	}
	title := r.PostForm.Get("title")
	body := r.PostForm.Get("body")
	errors := validateArticleFormData(title, body)

	if len(errors) != 0 {
		storeUrl := route.RouteName2URL("articles.store")
		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
		if err != nil {
			panic(err)
		}
		data := new(ArticlesFormData)
		data.URL = storeUrl
		data.Errors = errors
		data.Body = body
		tmpl.Execute(w, data)
		return
	}
	_article := &article.Article{
		Title: title,
		Body:  body,
	}
	_article.Create()
	if _article.ID > 0 {
		fmt.Fprintf(w, "数据插入成功id为%d", _article.ID)
	} else {
		fmt.Fprintf(w, "文章创建失败")
	}
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

func validateArticleFormData(title string, body string) map[string]string {
	errors := make(map[string]string)
	// 验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需介于 3-40"
	}

	// 验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需大于或等于 10 个字节"
	}

	return errors
}
