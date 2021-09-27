package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"go_blog/app/http/controllers"
	"net/http"
)

func RegisterWebRouter(router *mux.Router) *mux.Router {
	ArticleController := new(controllers.ArticlesController)
	router.HandleFunc("/", ArticleController.Index).Methods("GET").Name("articles.home")
	router.HandleFunc("/articles/{id:[0-9]+}", ArticleController.Show).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", ArticleController.Store).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles", ArticleController.Create).Methods("GET").Name("articles.create")
	router.HandleFunc("/articles/{id:[0-9]+}/edit", ArticleController.Edit).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", ArticleController.Update).Methods("POST").Name("articles.update")
	router.HandleFunc("/articles/{id:[0-9]+}/delete", ArticleController.Delete).Methods("GET").Name("articles.delete")
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	router.PathPrefix("/css/").Handler(http.FileServer(http.Dir("./public")))
	router.PathPrefix("/js/").Handler(http.FileServer(http.Dir("./public")))
	return router
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}