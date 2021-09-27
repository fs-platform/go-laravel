package main

//GO111MODULE=on go get -u github.com/cosmtrek/air
import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"go_blog/app/http/middlewares"
	"go_blog/bootstrap"
	"go_blog/pkg/databases"
	"net/http"
)

var db *sql.DB
var router *mux.Router

func main() {
	databases.Initialize()
	router = bootstrap.SetupRoute()
	bootstrap.SetupDB()
	db = databases.DB
	//自定义404
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	http.ListenAndServe(":3000", middlewares.RemoveSlash(router))
}


func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}