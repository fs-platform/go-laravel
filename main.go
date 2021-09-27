package main

//GO111MODULE=on go get -u github.com/cosmtrek/air
import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"go_blog/bootstrap"
	"go_blog/pkg/databases"
	"go_blog/pkg/logger"
	"go_blog/pkg/route"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"
)

var db *sql.DB
var router *mux.Router

func main() {
	databases.Initialize()
	router = bootstrap.SetupRoute()
	bootstrap.SetupDB()
	db = databases.DB
	router.HandleFunc("/", homeHandler).Methods("GET").Name("articles.home")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesEditHandler).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesUpdateHandler).Methods("POST").Name("articles.update")
	middlewares := []mux.MiddlewareFunc{
		setHeaderMiddleware,
	}
	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("GET").Name("articles.delete")
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")
	// 中间件：强制内容类型为 HTML
	router.Use(middlewares...)
	//自定义404
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	http.ListenAndServe(":3000", removeSlash(router))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		tmp      *template.Template
		rows     *sql.Rows
		articles []Article
	)
	tmp, err = template.ParseFiles("resources/views/articles/index.gohtml")
	logger.LogError(err)
	query := "SELECT * FROM articles"
	rows, err = db.Query(query)
	defer rows.Close()
	logger.LogError(err)
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.Title, &article.Body)
		logger.LogError(err)
		articles = append(articles, article)
	}
	err = rows.Err()
	logger.LogError(err)
	err = tmp.Execute(w, articles)
	logger.LogError(err)
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
	storeUrl, _ := router.GetRoute("articles.store").URL()
	tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}
	data := new(ArticlesFormat)
	data.URL = storeUrl
	err = tmpl.Execute(w, data)
	logger.LogError(err)
}

type Article struct {
	ID    int
	Body  string
	Title string
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	article, err := getArticleByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
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
	id := route.GetRouteVariable("id", r)
	editUrl, _ := router.GetRoute("articles.update").URL("id", id)
	article, _ := getArticleByID(id)
	data := ArticlesFormat{
		Title:  article.Title,
		Body:   article.Body,
		Errors: map[string]string{},
		URL:    editUrl,
	}
	err = tmpl.Execute(w, data)
	logger.LogError(err)
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

func getArticleByID(id string) (Article, error) {
	article := Article{}
	query := "SELECT * FROM articles WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	return article, err
}

func articlesUpdateHandler(w http.ResponseWriter, r *http.Request) {
	body := r.PostFormValue("body")
	title := r.PostFormValue("title")
	errors := validate(title, body)
	id := route.GetRouteVariable("id", r)
	if len(errors) != 0 {
		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		storeUrl, _ := router.GetRoute("articles.edit").URL("id", id)
		logger.LogError(err)
		articleInfo := ArticlesFormat{
			Body:   body,
			Title:  title,
			Errors: errors,
			URL:    storeUrl,
		}
		err = tmpl.Execute(w, articleInfo)
		logger.LogError(err)
		return
	}
	query := "UPDATE articles SET title=?,body=? WHERE id=?"
	result, err := db.Exec(query, title, body, id)
	if err != nil {
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 服务错误")
	}
	if n, _ := result.RowsAffected(); n > 0 {
		showURL, _ := router.Get("articles.show").URL("id", id)
		http.Redirect(w, r, showURL.String(), http.StatusFound)
		return
	} else {
		fmt.Fprint(w, "没有任何修改")
	}
}

func (a Article) Link(name string) string {
	showUrl, err := router.Get(name).URL("id", strconv.Itoa(a.ID))
	if err != nil {
		logger.LogError(err)
		return ""
	}
	return showUrl.String()
}

// Delete 方法用以从数据库中删除单条记录
func (a Article) Delete() (rowsAffected int64, err error) {
	rs, err := db.Exec("DELETE FROM articles WHERE id = " + strconv.Itoa(a.ID))

	if err != nil {
		return 0, err
	}

	// √ 删除成功，跳转到文章详情页
	if n, _ := rs.RowsAffected(); n > 0 {
		return n, nil
	}

	return 0, nil
}

func articlesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	articleInfo, err := getArticleByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Println("文章未找到")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("服务内部错误")
		}
	} else {
		rows, err := articleInfo.Delete()
		if err != nil {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("服务内部错误")
		} else {
			if rows > 0 {
				indexURL, _ := router.Get("articles.home").URL()
				http.Redirect(w, r, indexURL.String(), http.StatusFound)
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 文章未找到")
			}
		}
	}
}
