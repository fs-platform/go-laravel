package main

//GO111MODULE=on go get -u github.com/cosmtrek/air
import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"
)

var router *mux.Router = mux.NewRouter()

func main() {
	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	middlewares := []mux.MiddlewareFunc{
		setHeaderMiddleware,
	}
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")
	// 中间件：强制内容类型为 HTML
	router.Use(middlewares...)
	//自定义404
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	http.ListenAndServe(":3000", removeSlash(router))
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog！</h1>")
}
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">Aron.Yao@feisu.com</a>")
}
func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		// 解析错误，这里应该有错误处理
		fmt.Fprint(w, "请提供正确的数据！")
		return
	}
	title := r.PostForm.Get("title")
	body := r.PostForm.Get("body")
	errors := make(map[string]string)
	if body == "" {
		errors["body"] = "body不能为空"
	} else if utf8.RuneCountInString(body) > 10 {
		errors["body"] = "body长度要大于10"
	}
	if title == "" {
		errors["title"] = "title不能为空"
	}
	if len(errors) != 0 {
		storeUrl, _ := router.Get("articles.store").URL()
		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
		if err != nil {
			panic(err)
		}
		data := new(ArticlesFormat)
		data.URL = storeUrl
		data.Errors = errors
		data.Body = body
		tmpl.Execute(w, data)
		return
	}

	fmt.Fprintf(w, "POST PostForm: %v <br>", r.PostForm)
	fmt.Fprintf(w, "POST Form: %v <br>", r.Form)
	fmt.Fprintf(w, "title 的值为: %v", title)
}
func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "访问文章列表")
}
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}

type ArticlesFormat struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
	storeUrl, _ := router.Get("articles.store").URL()
	tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}
	data := new(ArticlesFormat)
	data.URL = storeUrl
	tmpl.Execute(w, data)
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprint(w, "文章 ID："+id)
}

func setHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. 设置标头
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// 2. 继续处理请求
		next.ServeHTTP(w, r)
	})
}

func removeSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/" {
			request.URL.Path = strings.TrimSuffix(request.URL.Path, "/")
		}
		next.ServeHTTP(writer, request)
	})
}
