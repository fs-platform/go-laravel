package bootstrap

import (
	"github.com/gorilla/mux"
)

// SetupRoute 路由初始化
func SetupRoute() *mux.Router {
	router := mux.NewRouter()
	//routes.RegisterWebRouter(router)
	return router
}
