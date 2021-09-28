package controllers

import (
	"database/sql"
	"fmt"
	"go_blog/app/models/article"
	"go_blog/pkg/logger"
	"go_blog/pkg/route"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
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
	articles, err := article.GetAll()
	fmt.Println(articles)
	if err != nil {
		// 数据库错误
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 服务器内部错误")
	} else {
		viewDir := "resources/views"
		files, err := filepath.Glob(viewDir + "/layouts/*.gohtml")
		newFiles := append(files, viewDir+"/articles/index.gohtml")
		tmp, err = template.ParseFiles(newFiles...)
		logger.LogError(err)
		err = tmp.ExecuteTemplate(w, "app", articles)
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
	body := r.PostFormValue("body")
	title := r.PostFormValue("title")
	errors := validateArticleFormData(title, body)
	id := route.GetRouteVariable("id", r)
	if len(errors) != 0 {
		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		storeUrl := route.RouteName2URL("articles.update", "id", id)
		articleInfo := ArticlesFormData{
			Body:   body,
			Title:  title,
			Errors: errors,
			URL:    storeUrl,
		}
		err = tmpl.Execute(w, articleInfo)
		logger.LogError(err)
		return
	}
	idInt, _ := strconv.Atoi(id)
	_article := &article.Article{
		Title: title,
		Body:  body,
		ID:    idInt,
	}
	affect, err := _article.Update()
	if err != nil {
		fmt.Fprint(w, "文章更新失败")
		return
	}
	if affect > 0 {
		showURL := route.RouteName2URL("articles.show", "id", id)
		http.Redirect(w, r, showURL, http.StatusFound)
		return
	} else {
		fmt.Fprint(w, "没有任何修改")
	}
}

func (*ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	articleInfo, err := article.Get(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Println("文章未找到")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
		}
	} else {
		idInt, _ := strconv.Atoi(id)
		_article := &article.Article{
			ID:    idInt,
			Body:  articleInfo.Body,
			Title: articleInfo.Title,
		}
		affect, err := _article.Delete()
		if err != nil {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("服务内部错误")
		} else {
			if affect > 0 {
				indexURL := route.RouteName2URL("articles.home")
				http.Redirect(w, r, indexURL, http.StatusFound)
				return
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 文章未找到")
			}
		}
	}
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
	fmt.Println(errors)
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
	tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
	if err != nil {
		log.Fatal(err)
	}
	id := route.GetRouteVariable("id", r)
	editUrl := route.RouteName2URL("articles.update", "id", id)
	article, _ := article.Get(id)
	data := ArticlesFormData{
		Title:  article.Title,
		Body:   article.Body,
		Errors: map[string]string{},
		URL:    editUrl,
	}
	err = tmpl.Execute(w, data)
	logger.LogError(err)
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
