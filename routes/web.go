package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"go_blog/app/http/controllers"
	"go_blog/app/http/middlewares"
	"net/http"
)

func RegisterWebRouter(router *mux.Router) *mux.Router {
	ArticleController := new(controllers.ArticlesController)
	RegisterController := new(controllers.RegisterController)
	LoginController := new(controllers.LoginController)
	router.HandleFunc("/", ArticleController.Index).Methods("GET").Name("articles.home")
	router.HandleFunc("/articles/{id:[0-9]+}", ArticleController.Show).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", ArticleController.Store).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles", ArticleController.Create).Methods("GET").Name("articles.create")
	router.HandleFunc("/articles/{id:[0-9]+}/edit", ArticleController.Edit).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", ArticleController.Update).Methods("POST").Name("articles.update")
	router.HandleFunc("/articles/{id:[0-9]+}/delete", ArticleController.Delete).Methods("GET").Name("articles.delete")
	router.HandleFunc("/register", RegisterController.Register).Methods("GET").Name("auth.register")
	router.HandleFunc("/doRegister", RegisterController.DoRegister).Methods("POST").Name("auth.doRegister")
	router.HandleFunc("/login", LoginController.Login).Methods("GET").Name("auth.login")
	router.HandleFunc("/doLogin", LoginController.DoLogin).Methods("POST").Name("auth.doLogin")
	router.HandleFunc("/loginOut", LoginController.LoginOut).Methods("get").Name("auth.loginOut")
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	router.Use(middlewares.StartSession)
	router.PathPrefix("/css/").Handler(http.FileServer(http.Dir("./public")))
	router.PathPrefix("/js/").Handler(http.FileServer(http.Dir("./public")))
	return router
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}
