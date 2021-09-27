package main

//GO111MODULE=on go get -u github.com/cosmtrek/air
import (
	"go_blog/app/http/middlewares"
	"go_blog/bootstrap"
	"net/http"
)

func main() {
	router := bootstrap.SetupRoute()
	bootstrap.SetupDB()
	//自定义404
	http.ListenAndServe(":3000", middlewares.RemoveSlash(router))
}
