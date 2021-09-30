package main

//GO111MODULE=on go get -u github.com/cosmtrek/air
import (
	"go_blog/app/http/middlewares"
	"go_blog/bootstrap"
	"go_blog/config"
	"net/http"
)

func init() {
	// 初始化配置信息
	config.Initialize()
}

func main() {
	router := bootstrap.SetupRoute()
	bootstrap.SetupDB()
	//自定义404
	http.ListenAndServe(":3000", middlewares.RemoveSlash(router))
}
