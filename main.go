package main

//GO111MODULE=on go get -u github.com/cosmtrek/air
import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"go_blog/bootstrap"
	"go_blog/pkg/databases"
	"net/http"
	"strings"
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