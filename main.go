package main

//GO111MODULE=on go get -u github.com/cosmtrek/air
import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"
)

var db *sql.DB
var router *mux.Router = mux.NewRouter()

func main() {
	initDB()
	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesEditHandler).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("POST").Name("articles.update")
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
	errors := validate(title, body)

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
	id, err := saveArticleToDB(title, body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "数据更新成功id为%d", id)
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
	err = tmpl.Execute(w, data)
	checkError(err)
}

type Article struct {
	ID    int
	Body  string
	Title string
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	id := getRouteVariable("id", r)
	article, err := getArticleByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	}
	tmpl, err := template.ParseFiles("resources/views/articles/show.gohtml")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, article)
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

func initDB() {
	var err error
	config := mysql.Config{
		User:                 "root",
		Passwd:               "root",
		Addr:                 "127.0.0.1:8889",
		Net:                  "tcp",
		DBName:               "go_blog",
		AllowNativePasswords: true,
	}
	db, err = sql.Open("mysql", config.FormatDSN())
	checkError(err)
	// 设置最大连接数
	db.SetMaxOpenConns(25)
	// 设置最大空闲连接数
	db.SetMaxIdleConns(25)
	// 设置每个链接的过期时间
	db.SetConnMaxLifetime(5 * time.Minute)
	err = db.Ping()
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func saveArticleToDB(title string, body string) (int64, error) {
	var (
		id   int64
		err  error
		rs   sql.Result
		stmt *sql.Stmt
	)
	stmt, err = db.Prepare("INSERT INTO articles (title, body) VALUES(?,?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	rs, err = stmt.Exec(title, body)
	// 4. 插入成功的话，会返回自增 ID
	if id, err = rs.LastInsertId(); id > 0 {
		return id, nil
	}
	return 0, err
}

func articlesEditHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
	if err != nil {
		log.Fatal(err)
	}
	id := getRouteVariable("id", r)
	editUrl, _ := router.GetRoute("articles.update").URL("id", id)
	article, _ := getArticleByID(id)
	data := ArticlesFormat{
		Title:  article.Title,
		Body:   article.Body,
		Errors: map[string]string{},
		URL:    editUrl,
	}
	err = tmpl.Execute(w, data)
	checkError(err)
}

func validate(title string, body string) map[string]string {
	errors := make(map[string]string)
	if body == "" {
		errors["body"] = "body不能为空"
	} else if utf8.RuneCountInString(body) > 10 {
		errors["body"] = "body长度要大于10"
	}
	if title == "" {
		errors["title"] = "title不能为空"
	}
	return errors
}

func getRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}

func getArticleByID(id string) (Article, error) {
	article := Article{}
	query := "SELECT * FROM articles WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	return article, err
}

func articlesUpdateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "更新文章")
}
