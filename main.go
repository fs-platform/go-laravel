package main

//GO111MODULE=on go get -u github.com/cosmtrek/air
import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
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
	}
	if title == "" {
		errors["title"] = "title不能为空"
	}
	if len(errors) != 0{
		fmt.Fprintf(w, "POST PostForm: %v <br>", errors)
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

func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <title>创建文章 —— 我的技术博客</title>
</head>
<body>
    <form action="%s?data=aron" method="post">
        <p><input type="text" name="title"></p>
        <p><textarea name="body" cols="30" rows="10"></textarea></p>
        <p><button type="submit">提交</button></p>
    </form>
</body>
</html>`
	storeUrl, _ := router.Get("articles.store").URL()
	fmt.Fprintf(w, html, storeUrl)
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
