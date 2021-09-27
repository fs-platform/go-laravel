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
	middlewares := []mux.MiddlewareFunc{
		setHeaderMiddleware,
	}
	// 中间件：强制内容类型为 HTML
	router.Use(middlewares...)
	//自定义404
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	http.ListenAndServe(":3000", removeSlash(router))
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


type Article struct {
	ID    int
	Body  string
	Title string
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
