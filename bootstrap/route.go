package bootstrap

import (
	"github.com/gorilla/mux"
	"go_blog/pkg/route"
	"go_blog/routes"
)

// SetupRoute 路由初始化
func SetupRoute() *mux.Router {
	router := mux.NewRouter()
	routes.RegisterWebRouter(router)
	route.SetRoute(router)
	return router
}
